package main

import (
	"context"
	"encoding/hex"
	"log"
	"strings"
	"time"

	sova_sdk_go "github.com/sova-network/sova-sdk-go"
	"github.com/sova-network/sova-sdk-go/generated"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const JETTON_TRANSFER_OPCODE = 0xf8a7ea5

func main() {
	// Create a new liteclient ton connection pool
	tonClient := liteclient.NewConnectionPool()

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	err = tonClient.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}
	// api client with full proof checks
	api := ton.NewAPIClient(tonClient, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)
	// bound all requests to single ton node
	ctx := tonClient.StickyContext(context.Background())

	// seed words of account
	// You should replace the words with your own seed words
	words := strings.Split("<replace with your seed words>", " ")

	// You should choose the right version of the wallet
	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.WalletAddress())

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!")

	// Balance should be more than 0.05 TON
	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	log.Println("wallet balance:", balance)

	// Create a new sova client
	sovaClient := sova_sdk_go.NewTestnetClient()

	// Create a new searcher
	searcher, err := sovaClient.Searcher(ctx)
	if err != nil {
		panic(err)
	}

	// Get tip addresses
	tipAddresses, err := searcher.GetTipAddresses(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("tipAddresses:", tipAddresses.GetAddress())

	// Subscribe to bundle results
	// Create a channel to receive bundle results
	resultCh := make(chan *types.BundleResult, 1)
	err = searcher.SubscribeBundleResult(ctx, func(br *types.BundleResult) {
		log.Println("receive bundle result")
		resultCh <- br
	})
	if err != nil {
		panic(err)
	}

	// Create a channel to receive mempool packets
	mempoolCh := make(chan *types.MempoolPacket, 1)

	// Set shard to subscribe
	shard, _ := hex.DecodeString("A000000000000000")

	log.Printf("jetton transfer opcode: %v\n", JETTON_TRANSFER_OPCODE)
	log.Println("Subscribe to i")
	err = searcher.SubscribeByInternalMsgBodyOpcode(ctx, 0, shard, JETTON_TRANSFER_OPCODE, func(mp *types.MempoolPacket) {
		log.Println("receive opcode packet")
		log.Printf("shard: %v\n", shard)

		mempoolCh <- mp
	})
	if err != nil {
		panic(err)
	}

	// Uncomment the following code to subscribe by address or opcode
	// Subscribe to opcode
	/*

				// Subscribe to workchain and shard
				err = searcher.SubscribeByWorkchainShard(ctx, 0, shard, func(mp *types.MempoolPacket) {
					log.Println("receive shard packet")
					mempoolCh <- mp
				})
				if err != nil {
					panic(err)
				}

			    // "0:hex_adddress1", "0:hex_address2" // 0 - workchain id, hex_address - hex encoded address
		    	_ = searcher.SubscribeByAddress(ctx, []string{"hex_address", "hex_address"}, func(mp *types.MempoolPacket) {
		    		log.Println("receive address packet: %v", mp)
		    	})

		    	// Subscribe to workchain
		    	err = searcher.SubscribeByWorkchain(ctx, 0, func(mp *types.MempoolPacket) {
		            log.Println("receive wc packet")
		            mempoolCh <- mp
		        })
		        if err != nil {
		            panic(err)
		        }
	*/

	var auctionID string

	// Wait for mempool packets and result
	for {
		select {
		// Receive mempool packets
		case msg := <-mempoolCh:
			log.Printf("got mempool packets: %v\n", len(msg.ExternalMessages))

			for _, extMsg := range msg.ExternalMessages {
				log.Printf("ext msg.GasSpend: %v\n", extMsg.GasSpent)
				log.Printf("ext msg.Shard: %v\n", hex.EncodeToString(extMsg.Shard))
				log.Printf("ext msg.StdSmcAddress: %v\n", hex.EncodeToString(extMsg.StdSmcAddress))

				// Get the destination address
				destAddr := address.MustParseAddr(tipAddresses.Address[0])
				log.Println("destAddr:", destAddr.StringRaw())

				log.Println("Create front transfer message")

				// Create front transfer message
				// Jetton master contract address kQCPu1TQE_cnYOYLLqFW7SNO8a8QVMgncM1lvV8_8bda-oqH
				token := jetton.NewJettonMasterClient(api, address.MustParseAddr("kQCPu1TQE_cnYOYLLqFW7SNO8a8QVMgncM1lvV8_8bda-oqH"))

				// find our jetton wallet
				tokenWallet, err := token.GetJettonWallet(ctx, w.WalletAddress())
				if err != nil {
					log.Fatal(err)
				}

				tokenBalance, err := tokenWallet.GetBalance(ctx)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("our jetton balance:", tokenBalance.String())

				amountTokens := tlb.MustFromDecimal("0.1", 9)

				comment, err := wallet.CreateCommentCell("Hello from sova sdk-go")
				if err != nil {
					log.Fatal(err)
				}

				// address of receiver's wallet (not token wallet, just usual)
				to := address.MustParseAddr("0QBlY0N2N51O_R1W5dC9giWhFgT6EETUnaMOS5Sqc10DEkvZ")
				transferPayload, err := tokenWallet.BuildTransferPayloadV2(to, to, amountTokens, tlb.ZeroCoins, comment, nil)
				if err != nil {
					log.Fatal(err)
				}

				// your TON balance must be > 0.05 to send
				msg := wallet.SimpleMessage(tokenWallet.Address(), tlb.MustFromTON("0.05"), transferPayload)

				// Build transfer message from your wallet to one of the tip addresses
				bounce := true
				tipMsg, err := w.BuildTransfer(destAddr, tlb.MustFromTON("0.03"), bounce, "Hello from sova sdk-go")
				if err != nil {
					log.Fatalln("Transfer err:", err.Error())
					continue
				}

				// Build external message
				tlbExtMsg, err := w.BuildExternalMessageForMany(ctx, []*wallet.Message{
					tipMsg,
					msg,
				})
				if err != nil {
					log.Fatalln("BuildExternalMessage err:", err.Error())
					continue
				}

				// Convert external message to cell
				cell, err := tlb.ToCell(tlbExtMsg)
				if err != nil {
					log.Fatalln("ToCell err:", err.Error())
					continue
				}

				// Serialize cell to BOC
				frontTx := cell.ToBOCWithFlags(false)
				log.Println("frontTx:", hex.EncodeToString(frontTx))

				// Send bundle
				res, err := searcher.SendBundle(ctx, &types.Bundle{
					Message: []*types.ExternalMessage{
						{
							Data: frontTx, // add frontTx
						},
						{
							Data: []byte(extMsg.Data), // add original message
						},
					},
					ExpirationNs: nil,
				})
				if err != nil {
					log.Println("SendBundle err:", err.Error())
					continue
				}

				log.Printf("got bundle response: %v\n", res)
				auctionID = res.Id

				time.Sleep(5 * time.Second)
			}

		// Receive bundle result
		case res := <-resultCh:

			// Check if the result is for the auction we are interested in
			if res.Id != auctionID {
				continue
			}

			log.Println("bundle result: ", res)

			// Check bundle result type
			if winResult, ok := res.Result.(*types.BundleResult_Win); ok {
				log.Println("bundle result auctionID: ", winResult.Win.AuctionId)
				log.Println("bundle result winner: ", winResult.Win.EstimatedNanotonTip)
			} else if looseResult, ok := res.Result.(*types.BundleResult_Loose); ok {
				log.Println("bundle result auctionID: ", looseResult.Loose.AuctionId)
			} else if failureResult, ok := res.Result.(*types.BundleResult_Failure); ok {
				log.Println("bundle result failure: ", failureResult.Failure.Reason)
			} else if dropResult, ok := res.Result.(*types.BundleResult_Drop); ok {
				log.Println("bundle result drop: ", dropResult.Drop.Reason)
			} else {
				log.Println("Unknown bundle result type")
			}
		}
	}
}
