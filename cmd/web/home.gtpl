<html>
	<head>
	<title></title>
	</head>
	<body>
        <h1>Home</h1>
		Hello {{.Username}} you are following {{.Following}} users.
		<form action="/logout" method="post">
			<input type="submit" value="Logout">
		</form>
		<h3> Follow </h3>
		<form action="/followUser" method="post">
			User's username:<input type="text" name="username">
			<input type="submit" value="Follow User">
		</form>
		<h3> Unfollow </h3>
		<form action="/deleteFollowing" method="post">
			User's username:<input type="text" name="username">
			<input type="submit" value="Unfollow User">
		</form>
		<div style="width:100%; height:10%">
		<h3> Post </h3>
		<form action="/createPost" method="post">
			Post Title:<input type="text" name="title">
			Post Content:<input type="text" name="content">
			<input type="submit" value="Create Post">
		</form>
		<div style="width:100%; height:10%">
		<h3> Feed </h3>
		{{if .Posts}}
			{{range .Posts}}
				<h3>{{.title}}</h3>
				<p>{{.content}}</p>
				<p> by {{.author}}</p>
			{{end}}
		{{else}}
			<p>You are not following anyone yet</p>
		{{end}}
	</body>
</html>