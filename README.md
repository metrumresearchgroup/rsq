# RSQ
#### Run R scripts in the background and fetch their results later.

With RSQ, you can send long-running chunks of R code into background tasks that run on a separate, then fetch their results later when they are ready. This is particularly useful for Shiny applications that trigger resource intensive calculations and block the user from performing other tasks in the interim.

## Getting started

We recommend using [Homebrew](https://brew.sh/) to help get started. This guide will assume you have Homebrew installed.

RSQ is written in Go, and so you'll need to have Go installed.
```
brew install go
```

After Go is installed, simply run the following to get the latest version of RSQ:
```
brew tap metrumresearchgroup/homebrew-tap
brew install rsq
```

Once you've made it this far, open a command prompt in the directory you wish the queue to run from. Then, run:
```
rsq start
```

The queue will now be running, listening for requests on a localhost port. You can then submit jobs to the queue and monitor their progress from within an R environment using [RRSQ].

## Resetting the queue

RSQ runs from the directory that you ran "rsq start" from. In that directory, you will see a "badger/" directory. This folder stores data for the past and present jobs in the queue. To "reset" the queue, stop the command-prompt process (with Ctrl + C) and delete then "badger/" directory. Then, simply run "rsq start" again.


[rrsq][https://github.com/metrumresearchgroup/rrsq]
