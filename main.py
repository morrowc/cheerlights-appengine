#!/usr/bin/env python
#
# Copyright 2007 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
import datetime
import os
import wsgiref.handlers

from google.appengine.ext import db
from google.appengine.ext import webapp
from google.appengine.ext.webapp import template

class ColorSet(db.Model):
  """A quick DataModel for storage of colorset operations.

  Store the:
    source-ip
    color
    time (server time)
  """
  sip = db.StringProperty(required=False)
  color = db.StringProperty(required=False)
  date = db.DateTimeProperty(required=False)

class Report(webapp.RequestHandler):
  def get(self):
    sip = self.request.remote_addr
    # If this is not a user looking for all records, show only
    # their SIP's records.
    where = 'WHERE sip = \'%s\' ' % sip

    alimit = self.request.get('a')
    if alimit:
      where = ''
      sip = ''

    limit = self.request.get('limit')
    try:
      limit = int(limit)
    except (TypeError, ValueError):
      pass

    if not limit or (type(limit) != type(1)):
      limit = 10

    data_query = ('SELECT * '
                  '  FROM ColorSet '
                  ' %s'
                  ' ORDER BY date DESC '
                  ' LIMIT %s' % (where, limit))


    data = db.GqlQuery(data_query)
    results = data.fetch(limit)

    # Set the template dict/dict up.
    template_values = { 'data': results,}

    path = os.path.join(os.path.dirname(__file__), 'templates', 'report.html')
    self.response.out.write(template.render(path, template_values))


class MainPage(webapp.RequestHandler):
  def get(self):
    
    # Collect the color data.
    cs = ColorSet()
    cs.color = self.request.get('color')
    cs.sip = self.request.remote_addr
    cs.date = datetime.datetime.now()

    # Store the color data.
    if cs.color:
      cs.put() 
      self.response.out.write('ok\n')
    else:
      self.response.out.write('not-ok\n')


def main():
  application = webapp.WSGIApplication([
                                        ('/', MainPage),
                                        ('/report', Report),
                                       ], debug=True)
  wsgiref.handlers.CGIHandler().run(application)


if __name__ == '__main__':
  main()
