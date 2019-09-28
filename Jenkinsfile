def src_harmony = 'harmony'
def src_mcl = 'mcl'
def src_bls = 'bls'
def src_cli = 'go-sdk'

node {
  def cpu_count = 4
  def harmony_version = null
  def passed = false

  try {
    // clean stage, remove all files in workspace
    stage('clean') {
      cleanWs()
    }
    // checkout stage, checkout needed repository from github
    stage('checkout') {
      def checkout_stages = [
        'harmony': {
          dir (src_harmony) {
            git branch:"${params.HARMONY_BRANCH}", url:'https://github.com/harmony-one/harmony'

            def version = sh (
                script: 'git rev-list --count HEAD',
                returnStdout: true,
                ).trim()
            def commit = sh (
                script: 'git describe --always --long --dirty',
                returnStdout: true,
                ).trim()
            harmony_version = "v${version}-${commit}"
          }
        },
        'mcl': {
          dir(src_mcl) {
            git url:'https://github.com/harmony-one/mcl'
          }
        },
        'bls': {
          dir(src_bls) {
            git url:'https://github.com/harmony-one/bls'
          }
        },
        'cli': {
          dir(src_cli){
            git branch:"${params.CLI_BRANCH}", url:"https://github.com/harmony-one/go-sdk"
          }
        }
      ]
      parallel(checkout_stages)
    }
    // build stage, build everything needed.
    stage('build_mcl') {
      dir(src_mcl) {
        sh "make -j${cpu_count}"
      }
    }
    stage('build_bls') {
      dir(src_bls) {
        sh "make BLS_SWAP_G=1 -j${cpu_count}"
      }
    }
    stage('build_harmony') {
      dir(src_harmony) {
        sh 'export PATH=$PATH:/usr/local/go/bin; scripts/go_executable_build.sh'
      }
    }
    stage('test_cli') {
      def cli_test_stages = [
        'launch_localnet': {
          dir(src_harmony){
            sh "./test/kill_node.sh"
            sh "./test/deploy.sh -D ${params.LOCALNET_DURATION} ./test/configs/local-resharding.txt"
            sh "./test/kill_node.sh"
          }
        },
        'launch_cli_tests': {
          dir(src_cli){
            sh "./jenkinsTest.sh"
          }
        }
      ]
      cli_test_stages.failFast = true
      parallel(cli_test_stages)
    }
  }
  catch (exc) {
    currentBuild.result = 'FAILURE'
  }
  finally {
    stage('kill_localnet') {
      dir(src_harmony){
        sh "./test/kill_node.sh"
      }
    }
  }
}