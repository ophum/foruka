<!doctype html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<title>foruka</title>
</head>
<body>
<h1>foruka</h1>
<hr>
	@if (Route::has('login'))
		@auth
			<a href="{{ url('/home') }}">Home</a>
		@else
 			<a href="{{ route('login') }}">Login</a>

 			@if (Route::has('register'))
				<a href="{{ route('register') }}">Register</a>
			@endif
 		@endauth
	@endif
    </body>
</html>
