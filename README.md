# appgroup ![CI Build](https://github.com/brunograsselli/appgroup/actions/workflows/push.yml/badge.svg)

Go package to manage go routine groups. It blocks until the first go routine returns and then cancels the context. After the context is canceled, it blocks until all the go routines have returned or the shutdown timeout is reached.

It is a variation of `sync.errgroup`. The main differences are:
* It returns regardless of the go routine result (error or success)
* It controls shutdown timeout
