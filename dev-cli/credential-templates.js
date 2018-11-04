const {createContext, CryptoFactory} = require('sawtooth-sdk/signing');
const {Secp256k1PrivateKey, Secp256k1PublicKey} = require('sawtooth-sdk/signing/secp256k1');
const {createHash} = require('crypto');
const {protobuf} = require('sawtooth-sdk');
const request = require('request');
const colors = require('colors');
const proto = require('google-protobuf');
const any = require('google-protobuf/google/protobuf/any_pb.js');
const payloads_pb = require('./proto/payload_pb');
const ethUtil = require('ethereumjs-util');
const context = createContext('secp256k1');
const opn = require('opn');
const axios = require('axios');
const prettyjson = require('prettyjson');
const moment = require('moment');

// hard coded example private key
const pk = Buffer.from("e3ddee618d8a8864481e71021e42ed46c3ab410ab1ad7cdf0ff31f6d61739275", 'hex');
const priv = new Secp256k1PrivateKey(pk);
const signer = new CryptoFactory(context).newSigner(priv);

const CONFIG = {
    templates: {
        familyName: "credential-templates",
        familyVersion: "1.0",
        namespaces: {
            prefixes: {
                "templates": createHash('sha512').update("credential:templates").digest('hex').substring(0, 6),
            },
            templateAddress(owner, name, version) {
                const prefix = this.prefixes.templates;
                const o = createHash('sha512').update(owner).digest('hex').substring(0, 30);
                const n = createHash('sha512').update(name).digest('hex').substring(0, 30);
                const v = createHash('sha512').update(version).digest('hex').toLowerCase().substring(0, 4);
                return `${prefix}${o}${n}${v}`
            },
        }
    }
};

const newRPCRequest = (params, method) => {
    const payload = new payloads_pb.RPCRequest();
    payload.setMethod(method);
    payload.setParams(params);
    return payload;
};

const newTemplate = (owner) => {
    return {
        "name":"Hello World Template",
        "version": "v1",
        "owner": owner,
        "data": JSON.stringify({"hello": "world"}),
    };
};

const create = async (amount, recipient) => {
    const owner = ethUtil.pubToAddress(signer.getPublicKey().asBytes(), true).toString('hex');
    const template = newTemplate(owner);

    const axioscRequest = newRPCRequest(JSON.stringify(template), payloads_pb.Method.CREATE);
    const axioscRequestBytes = axioscRequest.serializeBinary();

    //compute state address
    const stateAddress = CONFIG.templates.namespaces.templateAddress(template.owner, template.name, template.version);

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf.TransactionHeader.encode({
        familyName: CONFIG.templates.familyName,
        familyVersion: CONFIG.templates.familyVersion,
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
        payloadSha512: createHash('sha512').update(axioscRequestBytes).digest('hex')
    }).finish();

    console.log(colors.yellow(`state addresses: ${stateAddress}`));
    return await submitTransaction(transactionHeaderBytes, axioscRequestBytes);
};

const deleteTemplates = async (addresses) => {
    const axioscRequest = newRPCRequest(JSON.stringify({addresses}), payloads_pb.Method.DELETE);
    const axioscRequestBytes = axioscRequest.serializeBinary();

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf.TransactionHeader.encode({
        familyName: CONFIG.templates.familyName,
        familyVersion: CONFIG.templates.familyVersion,
        inputs: addresses,
        outputs: addresses,
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
        payloadSha512: createHash('sha512').update(axioscRequestBytes).digest('hex')
    }).finish();

    console.log(colors.yellow(`state addresses: ${addresses}`));
    return await submitTransaction(transactionHeaderBytes, axioscRequestBytes);
};

const updateTemplate = async (amount, recipient) => {
    const owner = ethUtil.pubToAddress(signer.getPublicKey().asBytes(), true).toString('hex');
    const template = newTemplate(owner);

    const axioscRequest = newRPCRequest(JSON.stringify(template), payloads_pb.Method.CREATE);
    const axioscRequestBytes = axioscRequest.serializeBinary();

    //compute state address
    const stateAddress = CONFIG.templates.namespaces.templateAddress(template.owner, template.name, template.version);

    // do the sawtooth thang ;)
    const transactionHeaderBytes = protobuf.TransactionHeader.encode({
        familyName: CONFIG.templates.familyName,
        familyVersion: CONFIG.templates.familyVersion,
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
        payloadSha512: createHash('sha512').update(axioscRequestBytes).digest('hex')
    }).finish();

    console.log(colors.yellow(`state addresses: ${stateAddress}`));
    return await submitTransaction(transactionHeaderBytes, axioscRequestBytes);
};

const submitTransaction = async (transactionHeaderBytes, axioscRequestBytes) => {
    try {
        const signature = signer.sign(transactionHeaderBytes);

        const transaction = protobuf.Transaction.create({
            header: transactionHeaderBytes,
            headerSignature: signature,
            payload: axioscRequestBytes
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

        const reqConfig = {
            method: 'POST',
            url: 'http://127.0.0.1:8008/batches',
            data: batchListBytes,
            headers: {'Content-Type': 'application/octet-stream'}
        };

        const response = await  axios(reqConfig);

        const link = response.data.link;
        console.log(colors.green(`transaction submitted successfully`));
        console.log(colors.green(`status: ${link}`));
        opn(link);
        process.exit(0);
    } catch (e) {
        throw e;
    }
};

const queryState = async (address) => {
    try {
        const reqConfig = {
            method: 'GET',
            url: `http://127.0.0.1:8008/state?address=${address}`,
            headers: {'Content-Type': 'application/json'}
        };

        const response = await  axios(reqConfig);
        const data = response.data.data;
        data.forEach(entry => {
            const data = new Buffer(entry.data, 'base64').toString('ascii');
            const template = JSON.parse(data);

            const output = {
                'state-address': entry.address,
                'template': template,
                'created-at': moment(new Date(template["created_at"] * 1000)).format('L'),
            };
            console.log(prettyjson.render(output));
        });

        process.exit(0);
    } catch (e) {
        throw e;
    }
};

module.exports = {
    create,
    deleteTemplates,
    queryState
};
