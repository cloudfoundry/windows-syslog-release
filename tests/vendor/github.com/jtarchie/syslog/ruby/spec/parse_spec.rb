# frozen_string_literal: true
require 'spec_helper'
require 'benchmark'
require_relative '../parse'

RSpec.describe 'Parse syslog messages' do
  it 'parses sucessfully' do
    payload = %(<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8)
    log = Syslog.new.parse(payload)
    expect(log[:message]).to eq "'su root' failed for lonvick on /dev/pts/8"
  end

  it 'is fast' do
    payload = %(<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8)
    syslog = Syslog.new

    puts(Benchmark.bm do |x|
      x.report("parse") {
        (1..100).each do
          syslog.parse(payload)
        end
      }
    end)
  end
end
