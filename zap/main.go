package main

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Infow("testing", "A", "a", "B", "b")

	logger.Info("testing2", zap.String("A", "a"))
}
