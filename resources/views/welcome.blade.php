@extends('layouts.app')

@section('content')

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
@endsection