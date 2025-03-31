# gonuxt-starter

## Installing

### Step 1. Docker postgres image (optional)

*This step assumes you have docker already installed on your machine*</br>

I recommend outsourcing the database and skipping to step 3.</br>
Self hosting a database in a production environment is easier said than done, even if everything seems to be working properly :)


<h4>1. Pulling the latest postgres image for docker</h4>

```
docker pull postgres:latest
```

<h4>2. Running your image on a container</h4>

**Syntax**
```
docker run --name [CONTAINER_NAME] -p [INTERNAL PORT:EXTERNAL PORT] -e POSTGRES_USER=[DATABASE USER NAME] -e POSTGRES_PASSWORD=[DATABASE USER PASSWORD] -d postgres
```
**My defaults**
```
docker run --name gonuxt-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=gonuxtsecret -d postgres
```

<h4>3. Connecting with your database</h5>

<p>I recommend using psql, you would need to <a href=https://www.postgresql.org/download>install postgres</a> on your local machine.</p>

**Syntax**
```
psql [DATABASE DRIVER]://[DATABASE USER]:[DATABASE PASSWORD]@[IP ADDRESS]:[PORT]
```

**My defaults**
```
psql postgres://root:gonuxtsecret@localhost:5432
```

If your terminall connects to the database shell, it means you have completed the setup correctly.<br/>
You can exit out of the shell with the command <code>\q</code>

<h4>4. Creating the database for our project</h4>

<p>I created some Makefiles to make the setup a bit more intuitive.</p>
<p>You can find the commands ran inside the Makefiles at the root at the project.</p>
<p>In order to execute them, you will need <a href=https://www.gnu.org/software/make/> gnu make </a> </p>
<p>Based on your setup, you might want to modify the variables available inside <code>Makefile.variable</code></p>
<p>After you are done modifying the variables, you can run the command 
  
  ```
  make createdb
  ```
</p>
<p>This should create a database with the name set up inside the <code>Makefile.variable</code></p>

<h4>5. Creating migrations</h4>

I am using <a href=https://docs.sqlc.dev/en/latest/overview/install.html>sqlc</a> and <a href=https://pressly.github.io/goose/installation/>goose</a> for this project. Make sure to install them.
I created a folder called <code>db</code> with 3 other folders inside it. 
<ul>
  <li><code>db/schema</code> - we write the migrations we want to apply here</li>
  <li><code>db/queries</code> - we write the queries for our database here</li>
  <li><code>db/out</code> - the output folder where sqlc generates our go code</li>
</ul>
<p></p>You can find the syntax for writing migrations and queries on the official goose website. 
I created some files for our small use case. </p>
