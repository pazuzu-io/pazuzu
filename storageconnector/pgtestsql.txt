# simple sql querys to create a demo table for pazuzu in postgres

CREATE TABLE IF NOT EXISTS features (index int, name text, description text, author text, lastupdate timestamptz, dependencies text, snippet text);
INSERT INTO features VALUES (1, 'java', 'openjdk8', 'Olaf', '2016-11-25 11:36:25+01', '', 'sudo apt install openjdk-8-jre');
INSERT INTO features VALUES (2, 'python2', 'python2', 'Olaf', '2016-11-25 11:41:20+01', '', 'sudo apt install python');
INSERT INTO features VALUES (3, 'python3', 'python3', 'Olaf', '2016-11-25 11:41:20+01', '', 'sudo apt install python3');
INSERT INTO features VALUES (4, 'numpy', 'python2 package numpy', 'Olaf', '2016-11-25 11:43:50+01', 'python2', 'sudo apt install python-numpy');
INSERT INTO features VALUES (5, 'numpy3', 'python3 package numpy', 'Olaf', '2016-11-25 11:44:40+01', 'python3', 'sudo apt install python3-numpy');
INSERT INTO features VALUES (6, 'numpy-all', 'python2+python3 package numpy','Olaf', '2016-11-25 11:45:30+01    ', 'python2 python3', 'sudo apt install python-numpy && sudo apt install python3-numpy');

