# cheerlights-appengine
Appengine backend for cheerlights light driver.

This transitioned from apython/django based appengine site to golang.
The paired python script:
  https://github.com/morrowc/cheerlights-twitter

which connects to twitter to gather new #cheerlights commands, and
connects to the #cheerlights system over a serial/xbee interface.
