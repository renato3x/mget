package cli

import "flag"

type CmdArgs struct {
	Url       string
	AudioOnly bool
}

func Args() *CmdArgs {
  var args CmdArgs
  
  // Get user home directory and set fixed output directory
  // homeDir, _ := os.UserHomeDir()
  // args.Output = filepath.Join(homeDir, "mget-downloads")
  
  flag.BoolVar(&args.AudioOnly, "audio", false, "Download only the audio")
  flag.BoolVar(&args.AudioOnly, "a", false, "Download only the audio")
 
  flag.Parse()

  positionalArgs := flag.Args()
  if len(positionalArgs) > 0 {
    args.Url = positionalArgs[0]
  }

  return &args
}
