pipeline {
    agent any
    parameters {

        choice(name: 'OS', choices: ['linux', 'darwin', 'windows', 'all'], description: 'Pick OS')

    }
    stages {
        stage('Example') {
            steps {
                echo "Build for platform ${params.OS}"

                echo "Build for arch: ${params.ARCH}"

            }
        }
    }
}
