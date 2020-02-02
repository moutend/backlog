#!/usr/bin/env ruby

# Go's `ast` package is too much for this task, so I decided to use simple Ruby script.

functions = []

ARGF.each do |line|
  # Extract function statements.
  next if !line.start_with? "func (c *Client)"

  line = line.sub("func (c *Client) ", "")

  # Skip private methods
  next if /[a-z]/.match line[0]

  p_begin = line.index "("
  p_end = line.index ")"

  f_name = line[0 ... p_begin]
  f_args = line[(p_begin  + 1) .. (p_end - 1)]
  f_ret = line[(p_end + 2) ... -3]
  f_params = []

  f_args.split(",").each do |arg|
    f_params.push(arg.split(" ")[0])
  end

  next if f_name == "SetHTTPClient"

  functions.push({
    :name => f_name,
    :args => f_args,
    :ret => f_ret,
    :params => f_params,
  })
end

puts <<END
package backlog

import (
  . "github.com/moutend/go-backlog/pkg/types"
)
END

functions.each do |f|
  params = f[:params].join(", ")

  puts ""
  puts "func #{f[:name]}(#{f[:args]}) #{f[:ret]} {"
  puts "  return bc.#{f[:name]}(#{params})"
  puts "}"
end
