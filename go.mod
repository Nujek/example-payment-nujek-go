module github.com/Nujek/example-merchant-client

go 1.22

require github.com/Nujek/sdk-payment-nujek-go v0.1.0

// Local replace membuat contoh ini langsung runnable dari checkout payment-core.
// Hapus baris ini saat mengambil SDK dari proxy Git/module publik.
replace github.com/Nujek/sdk-payment-nujek-go => ../../sdk/go/merchantapi
