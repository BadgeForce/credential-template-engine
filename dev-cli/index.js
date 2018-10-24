const {createContext, CryptoFactory} = require('sawtooth-sdk/signing');
const {Secp256k1PrivateKey} = require('sawtooth-sdk/signing/secp256k1');
const {createHash} = require('crypto');
const {protobuf} = require('sawtooth-sdk');
const request = require('request');
const colors = require('colors');
const proto = require('google-protobuf');
const any = require('google-protobuf/google/protobuf/any_pb.js');
const payloads_pb = require('./proto/payload_pb');
const context = createContext('secp256k1');

// hard coded example private key
const pk = Buffer.from("e3ddee618d8a8864481e71021e42ed46c3ab410ab1ad7cdf0ff31f6d61739275", 'hex');
const priv = new Secp256k1PrivateKey(pk);
const signer = new CryptoFactory(context).newSigner(priv);

const CONFIG = {
    earnings: {
        familyName: "bf-credential-templates",
        familyVersion: "1.0",
        namespaces: {
            prefix: 'credential:templates',
            stateAddress(owner, name) {
                const prefix = createHash('sha512').update(this.prefix).digest('hex').substring(0, 6);
                const ownerPrefix = createHash('sha512').update(owner).digest('hex').substring(0, 32);
                const postfix = createHash('sha512').update(name).digest('hex').substring(0, 32);
                return `${prefix}${ownerPrefix}${postfix}`
            }
        }
    }
};

const newRPCRequest = (params, method) => {
    const payload = new payloads_pb.RPCRequest();
    payload.setMethod(method);
    payload.setParams(JSON.stringify(params));

    return payload;
};

const newTemplateReq = (owner, name, version, data) => {
    return {owner, name, version, data: JSON.stringify(data)};
};

const createTemplate = () => {
    //setup details
    const data = {
        position: "doctor",
        experience: "8 years",
        education: "PHD"
    };
    const name = "Test Template";
    const owner = signer.getPublicKey().asHex();
    const version = "v1";

    const template = newTemplateReq(owner, name, version, data);
    const rpcRequest = newRPCRequest(template, payloads_pb.Method.CREATE);
    const rpcRequestBytes = rpcRequest.serializeBinary();
    //compute state address
    const stateAddress = CONFIG.earnings.namespaces.stateAddress(owner, name);

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf.TransactionHeader.encode({
        familyName: CONFIG.earnings.familyName,
        familyVersion: CONFIG.earnings.familyVersion,
        inputs: [stateAddress],
        outputs: [stateAddress],
        signerPublicKey: signer.getPublicKey().asHex(),
        // In this example, we're signing the batch with the same private key,
        // but the batch can be signed by another party, in which case, the
        // public key will need to be associated with that key.
        batcherPublicKey: signer.getPublicKey().asHex(),
        // In this example, there are no dependencies.  This list should include
        // an previous transaction header signatures that must be applied for
        // this transaction to successfully commit.
        // For example,
        // dependencies: ['540a6803971d1880ec73a96cb97815a95d374cbad5d865925e5aa0432fcf1931539afe10310c122c5eaae15df61236079abbf4f258889359c4d175516934484a'],
        dependencies: [],
        payloadSha512: createHash('sha512').update(rpcRequestBytes).digest('hex')
    }).finish();

    submitTransaction(transactionHeaderBytes, rpcRequestBytes);
};

const submitTransaction = (transactionHeaderBytes, rpcRequestBytes) => {
    const signature = signer.sign(transactionHeaderBytes);

    const transaction = protobuf.Transaction.create({
        header: transactionHeaderBytes,
        headerSignature: signature,
        payload: rpcRequestBytes
    });

    const transactions = [transaction];

    const batchHeaderBytes = protobuf.BatchHeader.encode({
        signerPublicKey: signer.getPublicKey().asHex(),
        transactionIds: transactions.map((txn) => txn.headerSignature),
    }).finish();

    const signature1 = signer.sign(batchHeaderBytes);

    const batch = protobuf.Batch.create({
        header: batchHeaderBytes,
        headerSignature: signature1,
        transactions: transactions
    });

    const batchListBytes = protobuf.BatchList.encode({
        batches: [batch]
    }).finish();
    request.post({
        url: 'http://127.0.0.1:8008/batches',
        body: batchListBytes,
        headers: {'Content-Type': 'application/octet-stream'}
    }, (err, response) => {
        if (err) return console.log(err);
        const link = JSON.parse(response.body).link;

        console.log(colors.green(`transaction submitted successfully`));
        console.log(colors.green(`status: ${link}`));
        process.exit(0);
    });
};

const queryState = (address) => {
    request.get({
        url: `http://127.0.0.1:8008/state?address=${address}`,
        headers: {'Content-Type': 'application/json'}
    }, (err, response) => {
        if (err) {
            console.log(colors.red(err));
            process.exit(1);
        }

        const body = JSON.parse(response.body);
        const data = body.data;
        console.log(colors.green(`transaction submitted successfully`));
        console.log(`head: ${body.head}`);

        data.forEach(entry => {
           const bytes = new Uint8Array(Buffer.from(entry.data, 'base64'));
           const earning = new earnings_pb.Earning.deserializeBinary(bytes);

           console.log(`address: ${entry.address} `);
           console.log(colors.blue(`data: ${JSON.stringify(earning.toObject())}`));
        });

        process.exit(0);
    });
};

createTemplate();