all: test

venv:
	python3 -m venv venv
	source ./venv/bin/activate
	pip3 install -r requirements.txt

test:
	pytest -v ./ --disable-warnings
