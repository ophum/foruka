@extends('layouts.app')

@section('content')
<h1><a href="/">foruka</a></h1>
<h2>{{ __('Login') }}</h2>

<form method="POST" action="{{ route('login') }}">
@csrf
<p>
	<label for="email">{{ __('E-Mail Address') }}</label>

	<input id="email" type="email" class="{{ $errors->has('email') ? ' is-invalid' : '' }}" name="email" value="{{ old('email') }}" required autofocus>

@if ($errors->has('email'))
	<span class="invalid-feedback" role="alert">
		<strong>{{ $errors->first('email') }}</strong>
	</span>
@endif
</p>
<p>
	<label for="password">{{ __('Password') }}</label>

	<input id="password" type="password" class="{{ $errors->has('password') ? ' is-invalid' : '' }}" name="password" required>

	@if ($errors->has('password'))
		<span class="invalid-feedback" role="alert">
			<strong>{{ $errors->first('password') }}</strong>
		</span>
	@endif
</p>
<p>	
	<input type="checkbox" name="remember" id="remember" {{ old('remember') ? 'checked' : '' }}>
	<label for="remember">
		{{ __('Remember Me') }}
	</label>
</p>
<p>
	<button type="submit">
	{{ __('Login') }}
	</button>
	@if (Route::has('password.request'))
		<a href="{{ route('password.request') }}">
			{{ __('Forgot Your Password?') }}
		</a>
	@endif
</p>
</form>
@endsection
