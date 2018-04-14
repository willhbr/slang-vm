def process(file, outfile)
  contents = File.read(file)
  current_replacements = []
  section = ""
  in_generic = false

  output = ""

  contents.split("\n").each_with_index do |line, idx|
    if line.start_with? '__GENERIC__'
      raise "Already in generic section" if in_generic
      in_generic = true
      next
    end
    if line.start_with? '__DOTYPES__'
      types = line.scan(/\{(.*?)\}/)[0][0].split(/, ?/).map { |t| t.split /: ?/ }
      current_replacements << types
      next
    end
    if line.start_with? '__ENDGENERIC__'
      raise "Not in generic section" unless in_generic
      if current_replacements.empty?
        raise "No types defined in generic section"
      end
      current_replacements.each do |replacements|
        replacement = section
        replacements.each do |splat|
          type, special = splat
          special ||= ''
          replacement = replacement.gsub(type, special)
        end
        output << replacement
      end
      section = ''
      in_generic = false
      current_replacements = []
      next
    end

    if in_generic
      section << line + "\n"
    else
      output << line + "\n"
    end
  end
  File.write(outfile, output)
end

Dir[ARGV[0] + '/**/*.ggo'].each do |file|
  process(file, file.gsub(/\.ggo$/, '.go'))
end
