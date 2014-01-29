import json
import sys
from flask import Flask, url_for, render_template, request
import requests
from random import choice
import string

app = Flask(__name__)

examples = [
	'milf,teen',
	'hardcore,love',
	'mother,sister,brother,father',
	'secretary,boss',
	'blowjob,handjob,titjob,footjob'
	]

def format_query(query):
	s = ''
	for w in query.split(',')[:10]:
		if len(w) < 30:
			for ch in w:
				if ch in string.letters + string.digits:
					s += ch
			s += ','
	return s.strip(',')

def get_data(query):
	r = requests.get("http://localhost:8080/%s" % query)
	d = json.loads(r.text)
	res = [['Year'],['2008'],['2009'],['2010'],['2011'],['2012']]
	for k, v in d.iteritems():
		res[0].append(str(k))
		res[1].append(v['2008'])
		res[2].append(v['2009'])
		res[3].append(v['2010'])
		res[4].append(v['2011'])
		res[5].append(v['2012'])
	return str(res)

@app.route('/')
def index():
	try:
		q = format_query(request.args['q'])
	except:
		q = choice(examples)

	res = get_data(q)
	return render_template('index', q=q, res=res)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8081, debug=False)
