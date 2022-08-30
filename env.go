package go_chrome

import "os"

func init() {
	_ = os.Setenv("GOOGLE_API_KEY", "no")
	_ = os.Setenv("GOOGLE_DEFAULT_CLIENT_ID", "no")
	_ = os.Setenv("GOOGLE_DEFAULT_CLIENT_SECRET", "no")
}
