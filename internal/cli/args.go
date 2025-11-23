package cli

import "flag"

type CmdArgs struct {
	Url       string
	AudioOnly bool
	Version   bool
}

func Args() *CmdArgs {
  var args CmdArgs
  
  flag.BoolVar(&args.AudioOnly, "audio", false, "Download only the audio")
  flag.BoolVar(&args.AudioOnly, "a", false, "Download only the audio")
  flag.BoolVar(&args.Version, "version", false, "Show version information")
  flag.BoolVar(&args.Version, "v", false, "Show version information")
 
  flag.Parse()

  positionalArgs := flag.Args()
  if len(positionalArgs) > 0 {
    args.Url = positionalArgs[0]
  }

  return &args
}
