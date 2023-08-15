FROM ruby:slim
RUN gem install sqlite3 httparty
COPY main.rb /main.rb
CMD ["ruby", "main.rb"]