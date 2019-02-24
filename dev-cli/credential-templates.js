const {createContext, CryptoFactory} = require('sawtooth-sdk/signing');
const {Secp256k1PrivateKey, Secp256k1PublicKey} = require('sawtooth-sdk/signing/secp256k1');
const {createHash} = require('crypto');
const {protobuf} = require('sawtooth-sdk');
const request = require('request');
const colors = require('colors');
const proto = require('google-protobuf');
const any = require('google-protobuf/google/protobuf/any_pb.js');
const payloads_pb = require('./proto/payload_pb');
const template_pb = require('./proto/template_pb');
const transaction_receipts_pb = require('./proto/transaction_receipts_pb');
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
                console.log(prefix);
                const o = createHash('sha512').update(owner).digest('hex').substring(0, 30);
                const n = createHash('sha512').update(name).digest('hex').substring(0, 30);

                const vrsn = `${version.getMajor()}.${version.getMinor()}.${version.getPatch()}`;
                const v = createHash('sha512').update(vrsn).digest('hex').toLowerCase().substring(0, 4);
                return `${prefix}${o}${n}${v}`
            },
        }
    }
};

const createRPCRequest = (template) => {
    const create = new payloads_pb.Create()
    create.setParams(template);

    const payload = new payloads_pb.RPCRequest();
    payload.setCreate(create);

    return payload;
};

const getPOIHash = (templateData) => {
    const bytes = templateData.serializeBinary();
    return  createHash('md5').update(bytes).digest('hex');
} 

const newTemplate = (issuer, name, description, type, mjrNum, mnrNum, patchNum) => {
    const template = new template_pb.Template();
    const data = new template_pb.Data();
    const version = new template_pb.Version();
    const verification = new template_pb.Verification();

    version.setMajor(mjrNum);
    version.setMinor(mnrNum);
    version.setPatch(patchNum);

    data.setVersion(version);
    data.setCreatedAt(Date.now());
    data.setType(type);
    data.setDescription(description);
    data.setName(name);
    data.setIssuerPub(issuer.getPublicKey().asHex());

    const poiHash = getPOIHash(data);
    const sig =  issuer.sign(Buffer.from(poiHash, 'hex'));

    verification.setProofOfIntegrityHash(poiHash);
    verification.setSignature(sig);

    template.setData(data);
    template.setVerification(verification);

    return template;
};

const create = async (amount, recipient) => {
    const name = "First Test Template"
        description = "This is just a test baby"
        type = "TEST";
        mjrNum = "1"
        mnrNum = "0"
        patchNum = "0"

    const template = newTemplate(signer, name, description, type, mjrNum, mnrNum, patchNum);

    const rpcReq = createRPCRequest(template);
    const reqBytes = rpcReq.serializeBinary();

    //compute state address
    const stateAddress = CONFIG.templates.namespaces.templateAddress(template.getData().getIssuerPub(), template.getData().getName(), template.getData().getVersion());

    // // do the sawtooth thang ;)
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
        payloadSha512: createHash('sha512').update(reqBytes).digest('hex')
    }).finish();

    console.log(colors.yellow(`state addresses: ${stateAddress}`));
    return await submitTransaction(transactionHeaderBytes, reqBytes);
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
            const data = Buffer.from(entry.data, 'base64');
            const template_proto = template_pb.Template.deserializeBinary(data);
            const template = template_proto.toObject();

            const computedHash = getPOIHash(template_proto.getData())
            let integrity;
            if(computedHash === template.verification.proofOfIntegrityHash) {
                integrity = `VALID - (Matching Hashes) computed (${computedHash}) == expected (${template.verification.proofOfIntegrityHash})`;
            } else {
                integrity = `INVALID - (Mismatch Hashes) computed (${computedHash}) != expected (${template.verification.proofOfIntegrityHash})`;
            }
            
            const issuerPub = new Secp256k1PublicKey(Buffer.from(template.data.issuerPub, 'hex'));
            let sigValidation;

            if(context.verify(template.verification.signature, Buffer.from(computedHash, 'hex'), issuerPub)) {
                sigValidation = `VALID - template data ${computedHash} signed by ${issuerPub.asHex()} ownership verified`;
            } else {
                sigValidation = `INVALID - template data ${computedHash} not signed by ${issuerPub.asHex()} ownership could not be verified`;
            }

            const createdAt = template.data.createdAt;
            const version = template.data.version;
            const output = {
                'state-address': entry.address,
                'template': template,
                'created-at': moment(new Date(createdAt)).format('L'),
                'version': `${version.major}.${version.minor}.${version.patch}`,
                'data-integrity': integrity,
                'signature-validation': sigValidation
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
    //deleteTemplates,
    queryState
};
