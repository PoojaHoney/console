const jwt = require('jsonwebtoken');

console.log('Hello console');

const jwtPayload = {
  iss: 'console-automation@peaceful-harbor-416116.iam.gserviceaccount.com',
  sub: 'console-automation@peaceful-harbor-416116.iam.gserviceaccount.com',
  aud: 'https://oauth2.googleapis.com/token',
  exp: Math.floor(Date.now() / 1000) + 3600,
};

const private_key = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCm1WufS7sZaCKK\nynM4622gEHpCqlVbR5hbbGMNOprl27CJ5+bwQ3kSyDpGfAoE14iqTRXtDsqp1yl0\nmIOc2WgdfFNX/os2dB8uqmalbIjypZIIWxqLNBQ95TCMGPxCj5+rCTjyov77PoW+\nNNFkJuOFRMo9ELK+ki3SGsX9QPUju+X5Yy9aqDgNflbnX7bd9nJsOlMdYWsTb3VT\nYjbo3UofaSilE+cVycbd6cKbHOOsH3k3b+ew/GgzZ61f/teB2O2aacx6ftEfvZvK\nJNlgCGOMX335N1wAP+VNS1xqbC5mzwzfd4YzHHxmM3VEsLOoVe/K7ybm24vBDvS+\nw0uj3h1JAgMBAAECggEAEww56emRFpEDpJoBtsV2cjh/ZP4imbXeaM3Cr196EPVY\nvh4KiXMCr0jKEoMV98CN+3eqobK3I9YLhyLkn/NYhklMELdguJpgdwkAiDGQnzeV\nKTwRl0QebYSt2sj9gDH6tmHqrRU8DH5ycamBILCp+GzGtFskNUrmTV8+LLCObIJX\nuXe29ng7v11AnZaEFiQTH1OHWpKrJw8PR+jyU5duDwYbY/Kn74pufOd655u5BJS2\nd7jIBQdJVHsi2jl6S1kbM/ow0SQ0RjCQadslbGPqAXGWHrknJz6yXBZw/pTfkRx4\nf7EiP+iX/8+v98p7w1iRR5+o8SEGvesJNUK25VbxywKBgQDZydZHuTmuIbpdhRcM\nLQk5WstxjBpunRgxVbcpR8Qr1xJg4evbcdtHMO35CtY6R1+Ckd+GsYuJ8gIg0Ic7\ndRatP4uuoQ8oUENn8HZCSZGWR5CPOx3JrI5aieP9yyPnookU2ifqOJgQTimGSOiz\neX/cJiFhkPsNAk58WozPdZ3eqwKBgQDEGuej976xWSmAne4U2WUzESveVLO9+zNQ\nIzn4e/fS4NsIw/V3HUGM2niWHv3rkYtKZx3Oz7y9n8H7S8ylT3pVZKf/l/OSN18E\nQzXgZ8KWaa/dZywIqBo/c0zrh0CKB8HrWzc3oUKvTxTWhZfyTFnTy6NGnft1oRRb\nTNSucB3j2wKBgDQCqqS4TFkUfvBEj2t9+dzznPBB+DIAXD2z7ajzFZsuopn2eiAP\nAcHvonU+LEkAsIN7GLmO/vbzj6SMiC8f2EmJim8q0XbcLCzPVT6hWR5det1pzcRS\n7DAdDBCZCbsQqtILA0tBNrd2Ix6JnOD3nBxUEta9C+dgbKPv7/6/ZPM/AoGBAMLu\nn52MscBmieuwkS2mX4MgqYOqnLTXU81wBrrqt5CmNPQEniaWLUePD1jiW7NjDJub\n3rRqKQoxGMjsMvOMCmWB1cEq0VJhuhBpos97NKEBU3B4kvvT5at2uFpTKqNKTMff\n9wZURQ4wRN1tFHokzRHMFXQnrALkaHDM8YioJ/4LAoGBAIEkeNDS14jWmEwwmyQf\niWQkaSmgDMNQzvXPom+5yohPFeB6bvHrJ+3ubj2jSood2Uy1gotRDrDMzQJwQolh\nOA7gW9Y7cA+OQSWLotEFR5+2I5cKB6KCsoNol5GrZoCRvIbl6wcB/aD80wVe0/Ma\n/04r5ir87yHXNT0uvBO2Iwt/\n-----END PRIVATE KEY-----\n"

const jwtToken = jwt.sign(
  jwtPayload, private_key, 
 { algorithm: 'RS256' }
);

console.log(jwtToken);
