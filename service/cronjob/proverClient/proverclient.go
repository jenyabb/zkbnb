package main

import (
	"flag"
	"fmt"
	"github.com/zecrey-labs/zecrey-crypto/zecrey-legend/circuit/bn254/block"
	"github.com/zecrey-labs/zecrey-legend/common/util"
	"github.com/zecrey-labs/zecrey-legend/service/cronjob/proverClient/internal/config"
	"github.com/zecrey-labs/zecrey-legend/service/cronjob/proverClient/internal/logic"
	"github.com/zecrey-labs/zecrey-legend/service/cronjob/proverClient/internal/svc"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f",
	"./etc/proverClient.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	// srv := server.NewProverClientPingServer(ctx)

	cronJob := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DiscardLogger),
	))
	// init r1cs
	var circuit block.BlockConstraints
	fmt.Println("start compile circuit")
	r1csValue, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit, frontend.IgnoreUnconstrainedInputs())
	if err != nil {
		panic("r1cs init error")
	}
	fmt.Println("circuit constraints:", r1csValue.GetNbConstraints())
	fmt.Println("finish compile circuit")
	// read proving and verifying keys
	provingKey, err := util.LoadProvingKey(c.KeyPath.ProvingKeyPath)
	if err != nil {
		panic("provingKey loading error")
	}
	verifyingKey, err := util.LoadVerifyingKey(c.KeyPath.VerifyingKeyPath)
	if err != nil {
		panic("verifyingKey loading error")
	}

	_, err = cronJob.AddFunc("@every 10s", func() {
		// cron job for receiving cryptoBlock and handling
		err = logic.ProveBlock(ctx, r1csValue, provingKey, verifyingKey)
		if err != nil {
			logx.Error("Prove Error: ", err.Error())
		}
	})
	if err != nil {
		panic(err)
	}
	cronJob.Start()

	logx.Info("proverClient cronjob is starting......")
	select {}
}
