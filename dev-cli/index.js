const credentialTemplates = require('./credential-templates');
const cli = require('caporal');
const figlet = require('figlet');

cli
    .version('0.0.1')
    .command('create-template', 'create a new credential template')
    // .argument('<recipient>', 'ethereum address of the recipient')
    // .argument('<amount>', 'amount of props to issue in earning')
    .action(async (args, options, logger) => {
        logger.info(`submitting create template transaction`);
        try {
            await credentialTemplates.create();
        } catch (e) {
            logger.error(`error creating template: ${e}`)
        }
    })
    .command('delete-templates', 'delete some template(s) from the state')
    .argument('<addresses...>', 'state address for query', cli.LIST)
    .action(async (args, options, logger) => {
        try {
            await credentialTemplates.create(args.addresses);
        } catch (e) {
            logger.error(`error deleting template(s): ${e}`)
        }
    })
    .command('query-templates', 'get some template(s) from the state')
    .argument('<stateaddress>', 'state address for query')
    .action(async (args, options, logger) => {
        try {
            await credentialTemplates.queryState(args.stateaddress);
        } catch (e) {
            logger.error(`error making state query: ${e}`)
        }
    });

const banner = figlet.textSync('BadgeForce-CLI', {
    font: 'slant',
    horizontalLayout: 'fitted',
    verticalLayout: 'default'
});

cli.description(`\n\n${banner}`);
cli.parse(process.argv);