package utils

import (
	"flag"
	"os"
)

func HandleFlag() {
	flag.Func("a", "HTTP server address", func(aFlagValue string) error {
		return os.Setenv("RUN_ADDRESS", aFlagValue)
	})

	flag.Func("d", "Address of db connection", func(dFlagValue string) error {
		return os.Setenv("DATABASE_URI", dFlagValue)
	})

	flag.Func("r", "Address of accrual system", func(rFlagValue string) error {
		return os.Setenv("ACCRUAL_SYSTEM_ADDRESS", rFlagValue)
	})
}
