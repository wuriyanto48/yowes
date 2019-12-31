const { Binary } = require("binary-install");
const os = require("os");
const { join } = require("path");

const getPlatform = () => {
  const type = os.type();
  const arch = os.arch();

  if (type === "Windows_NT" && arch === "x64") {
    return "windows-amd64";
  }
  if (type === "Linux" && arch === "x64") {
    return "linux-amd64";
  }
  if (type === "Darwin" && arch === "x64") {
    return "darwin-amd64";
  }

  throw new Error(`Unsupported platform: ${type} ${arch}`);
};

const getBinary = () => {
  const platform = getPlatform();
  const version = require("./package.json").version;
  const url = `https://github.com/wuriyanto48/yowes/releases/download/v${version}/yowes-v${version}.${platform}.tar.gz`
  const installDirectory = join(os.homedir(), ".yowes");
  return new Binary(url, { name: "yowes", installDirectory });
};

const run = () => {
  const binary = getBinary();
  binary.run();
};

const install = () => {
  const binary = getBinary();
  binary.install();
};

const uninstall = () => {
  const binary = getBinary();
  binary.uninstall();
}

module.exports = {
  install,
  run,
  uninstall
};