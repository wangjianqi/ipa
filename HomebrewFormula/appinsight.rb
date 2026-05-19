class Appinsight < Formula
  desc "iOS App analysis tool for developers"
  homepage "https://github.com/wangjianqi/ipa"
  version "0.1.0"

  # 根据发布平台选择
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/wangjianqi/ipa/releases/download/v#{version}/appinsight-darwin-arm64"
    sha256 "YOUR_SHA256_HERE"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/wangjianqi/ipa/releases/download/v#{version}/appinsight-darwin-amd64"
    sha256 "YOUR_SHA256_HERE"
  end

  license "MIT"

  depends_on "majd/repo/ipatool" => :recommended

  def install
    bin.install "appinsight"
  end

  test do
    system "#{bin}/appinsight", "--version"
  end
end
