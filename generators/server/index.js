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
        this.log(this.destinationRoot())
        this.log(this.sourceRoot())
        const name = this.config.get('name')
        // copy ssl
        this.fs.copy(this.templatePath('certs/test/Makefile'), './certs/test/Makefile')
        // copy protocol
        this.fs.copy(this.templatePath('protocol/hello.pb.go'), './protocol/hello.pb.go')
        this.fs.copy(this.templatePath('protocol/hello.pb.gw.go'), './protocol/hello.pb.gw.go')
        this.fs.copy(this.templatePath('protocol/hello.proto'), './protocol/hello.proto')
        // copy server
        this.fs.copyTpl(
            this.templatePath('server/rpc.defs.go'),
            this.destinationPath('./server/rpc.defs.go'),
            { appname: name }
        ),
        this.fs.copyTpl(
            this.templatePath('server/serve_test.go'),
            this.destinationPath('./server/serve_test.go'),
            { appname: name }
        ),
        this.fs.copyTpl(
            this.templatePath('server/serve.go'),
            this.destinationPath('./server/serve.go'),
            { appname: name }
        ),
        // copy swagger
        this.fs.copy(this.templatePath('swagger/favicon-16x16.png'), './swagger/favicon-16x16.png'),
        this.fs.copy(this.templatePath('swagger/favicon-32x32.png'), './swagger/favicon-32x32.png'),
        this.fs.copy(this.templatePath('swagger/index.html'), './swagger/index.html'),
        this.fs.copy(this.templatePath('swagger/oauth2-redirect.html'), './swagger/oauth2-redirect.html'),
        this.fs.copy(this.templatePath('swagger/protocol.swagger.json'), './swagger/protocol.swagger.json'),
        this.fs.copy(this.templatePath('swagger/swagger-ui-bundle.js'), './swagger/swagger-ui-bundle.js'),
        this.fs.copy(this.templatePath('swagger/swagger-ui-bundle.js.map'), './swagger/swagger-ui-bundle.js.map'),
        this.fs.copy(this.templatePath('swagger/swagger-ui-standalone-preset.js'), './swagger/swagger-ui-standalone-preset.js'),
        this.fs.copy(this.templatePath('swagger/swagger-ui-standalone-preset.js.map'), './swagger/swagger-ui-standalone-preset.js.map'),
        this.fs.copy(this.templatePath('swagger/swagger-ui.css'), './swagger/swagger-ui.css'),
        this.fs.copy(this.templatePath('swagger/swagger-ui.css.map'), './swagger/swagger-ui.css.map'),
        this.fs.copy(this.templatePath('swagger/swagger-ui.js'), './swagger/swagger-ui.js'),
        this.fs.copy(this.templatePath('swagger/swagger-ui.js.map'), './swagger/swagger-ui.js.map'),
        // copy vendor
        this.fs.copy(this.templatePath('vendor/vendor.json'), './vendor/vendor.json'),
        // copy main dir
        this.fs.copy(this.templatePath('.gitignore'), './.gitignore'),
        this.fs.copy(this.templatePath('docker-test.sh'), './docker-test.sh'),
        this.fs.copy(this.templatePath('version'), './version'),
         this.fs.copyTpl(
            this.templatePath('README.md'),
            this.destinationPath('./README.md'),
            { appname: name }
        ),
        this.fs.copyTpl(
            this.templatePath('circle.yml'),
            this.destinationPath('./circle.yml'),
            { appname: name }
        ),
        this.fs.copyTpl(
            this.templatePath('Dockerfile'),
            this.destinationPath('./Dockerfile'),
            { appname: name }
        ),
        this.fs.copyTpl(
            this.templatePath('main.go'),
            this.destinationPath('./main.go'),
            { appname: name }
        )
    }

    end() {
        this.log(chalk.magenta('Thank you for using the generator. Happy coding!'))
    }
}
