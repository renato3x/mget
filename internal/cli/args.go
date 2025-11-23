package cli

import "flag"

type CmdArgs struct {
	Url       string
	Output    string
	AudioOnly bool
}

func Args() *CmdArgs {
  var args CmdArgs
  
  flag.StringVar(&args.Output, "output", "", "Output directory")
  flag.BoolVar(&args.AudioOnly, "audio", false, "Download only the audio")

  // shorthand args
  flag.StringVar(&args.Output, "o", "", "Output directory")
  flag.BoolVar(&args.AudioOnly, "a", false, "Download only the audio")
 
  flag.Parse()

  positionalArgs := flag.Args()
  if len(positionalArgs) > 0 {
    args.Url = positionalArgs[0]
  }

  return &args
}
