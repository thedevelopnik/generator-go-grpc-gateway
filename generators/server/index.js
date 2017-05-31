const util = require('util')
const path = require('path')
const Generator = require('yeoman-generator')
const chalk = require('chalk')
const yosay = require('yosay')

module.exports = class extends Generator {
    prompting() {
        const helloOne = chalk.magenta("'ALLO 'ALLO\n")
        const helloTwo = chalk.yellow('Welcome to the Branded Go Service Generator')
        this.log(yosay(`${helloOne}${helloTwo}`))
        this.log(chalk.red("Hey! You're doing this in $GOPATH/src/jaxf-fanatics.github.corp/apparel riiiiiight?'"))
        return this.prompt([{
            type: 'input',
            name: 'name',
            message: 'Your project name, lowercase and hyphenated please!',
            default: this.appname
        }]).then(answers => {
            this.config.set('name', answers.name)
            this.log('App name: ', answers.name)
            this.config.save()
        })
    }

    writing() {
        const name = this.config.get('name')
        // copy ssl
        this.copy('certs/test/Makefile', './certs/test/Makefile')
        // copy protocol
        this.copy('protocol/hello.pb.go', './protocol/hello.pb.go')
        this.copy('protocol/hello.pb.gw.go', './protocol/hello.pb.gw.go')
        this.copy('protocol/hello.proto', './protocol/hello.proto')
        // copy server
        this.copyTpl(
            this.templatePath('server/rpc.defs.go'),
            this.destinationPath('./server/rpc.defs.go'),
            { appname: name }
        ),
        this.copyTpl(
            this.templatePath('server/serve_test.go'),
            this.destinationPath('./server/serve_test.go'),
            { appname: name }
        ),
        this.copyTpl(
            this.templatePath('server/serve.go'),
            this.destinationPath('./server/serve.go'),
            { appname: name }
        ),
        // copy swagger
        this.copy('swagger/favicon-16x16.png', './swagger/favicon-16x16.png'),
        this.copy('swagger/favicon-32x32.png', './swagger/favicon-32x32.png'),
        this.copy('swagger/index.html', './swagger/index.html'),
        this.copy('swagger/oauth2-redirect.html', './swagger/oauth2-redirect.html'),
        this.copy('swagger/protocol.swagger.json', './swagger/protocol.swagger.json'),
        this.copy('swagger/swagger-ui-bundle.js', './swagger/swagger-ui-bundle.js'),
        this.copy('swagger/swagger-ui-bundle.js.map', './swagger/swagger-ui-bundle.js.map'),
        this.copy('swagger/swagger-ui-standalone-preset.js', './swagger/swagger-ui-standalone-preset.js'),
        this.copy('swagger/swagger-ui-standalone-preset.js.map', './swagger/swagger-ui-standalone-preset.js.map'),
        this.copy('swagger/swagger-ui.css', './swagger/swagger-ui.css'),
        this.copy('swagger/swagger-ui.css.map', './swagger/swagger-ui.css.map'),
        this.copy('swagger/swagger-ui.js', './swagger/swagger-ui.js'),
        this.copy('swagger/swagger-ui.js.map', './swagger/swagger-ui.js.map'),
        // copy vendor
        this.copy('vendor/vendor.json', './vendor/vendor.json'),
        // copy main dir
        this.copy('.gitignore', './.gitignore'),
        this.copy('docker-test.sh', './docker-test.sh'),
        this.copy('version', './version'),
         this.copyTpl(
            this.templatePath('README.md'),
            this.destinationPath('./README.md'),
            { appname: name }
        ),
        this.copyTpl(
            this.templatePath('circle.yml'),
            this.destinationPath('./circle.yml'),
            { appname: name }
        ),
        this.copyTpl(
            this.templatePath('Dockerfile'),
            this.destinationPath('./Dockerfile'),
            { appname: name }
        ),
        this.copyTpl(
            this.templatePath('main.go'),
            this.destinationPath('./main.go'),
            { appname: name }
        )
    }

    end() {
        this.log(chalk.magenta('Thank you for using the generator. Happy coding!'))
    }
}
