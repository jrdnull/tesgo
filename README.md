TesGo - Tesco API Wrapper
=========

TesGo is a partial wrapper for the [Tesco API](https://secure.techfortesco.com/tescoapiweb/secretlogin.aspx) written in Go.

TesGo supports the following commands:
* LOGIN
* PRODUCTSEARCH
* CHANGEBASKET
* LISTBASKET

Documentation for Beta 1 of the Tesco API is [available here](https://secure.techfortesco.com/tescoapiweb/Tesco%20Grocery%20API%20Beta%201%20Edition%20-%20REST%20Reference%20Guide%201.0.0.26.pdf).

Issues
------

The CHANGEBASKET command is not working as expected, a successful response is given but there is no change in the users Tesco basket. The issue has been reported and I am currently awaiting a response.


Installation
------------

Simply import the package into your project:

    import "github.com/jrdnull/tesgo"

and when you build that project with `go build`, it will be
downloaded and installed automagically.

Usage
-----

First you will require a [Tesco Shopping account](https://secure.tesco.com/register/), [Developer account](https://secure.techfortesco.com/tescoapiweb/secretlogin.aspx) and to have read and agreed to the [Terms and Conditions](http://www.techfortesco.com/tescoapiweb/terms.htm) of the Tesco API.

```go
    tesco = tesgo.New("<dev-key>", "<app key>")
	if err := tesco.Login("<email>", "<password>"); err != nil {
		fmt.Println(err)
		return
	}

	result, err := tesco.ProductSearch("milk", 1, false)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, product := range result.Products {
			fmt.Println(product.Name, product.Price)
		}
	}
```
Outputs:
```
	Tesco Pure British Semi Skimmed Milk 2L 1.48
	Tesco British Semi Skimmed Milk 2.272L 4Pints 1.29
	Tesco British Semi Skimmed Milk 1.136L 2Pints 0.89
	Tesco British Semi Skimmed Milk 3.480L 6Pints 1.89
...
```

License
-------

TesGo is distributed under the the BSD 2-Clause License:

> Copyright (C) 2013, Jordon Smith <jrd@mockra.net>
> All rights reserved.
>
> Redistribution and use in source and binary forms, with or without
> modification, are permitted provided that the following conditions
> are met:
> 1. Redistributions of source code must retain the above copyright
>    notice, this list of conditions and the following disclaimer.
> 2. Redistributions in binary form must reproduce the above copyright
>    notice, this list of conditions and the following disclaimer in the
>    documentation and/or other materials provided with the distribution.
>
> THIS SOFTWARE IS PROVIDED BY AUTHOR AND CONTRIBUTORS ``AS IS'' AND
> ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
> IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
> ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
> FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
> DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
> OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
> HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
> LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
> OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
> SUCH DAMAGE.
    