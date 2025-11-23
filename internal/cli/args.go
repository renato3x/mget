package cli

import "flag"

type CmdArgs struct {
	Url       string
	Output    string
	AudioOnly bool
}

func Args() *CmdArgs {
  var args CmdArgs
  
  // Default output directory is current directory
  defaultOutput := "."
  
  flag.StringVar(&args.Output, "output", defaultOutput, "Output directory")
  flag.BoolVar(&args.AudioOnly, "audio", false, "Download only the audio")

  // shorthand args
  flag.StringVar(&args.Output, "o", defaultOutput, "Output directory")
  flag.BoolVar(&args.AudioOnly, "a", false, "Download only the audio")
 
  flag.Parse()

  positionalArgs := flag.Args()
  if len(positionalArgs) > 0 {
    args.Url = positionalArgs[0]
  }

  return &args
}
