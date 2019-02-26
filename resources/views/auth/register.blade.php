@extends('layouts.app')

@section('content')
    <h1><a href="/">foruka</a></h1>
    <h2>{{ __('Register') }}</h2>

    <form method="POST" action="{{ route('register') }}">
    @csrf

    <p>
    <label for="name">{{ __('Name') }}</label>

    <input id="name" type="text" class="{{ $errors->has('name') ? ' is-invalid' : '' }}" name="name" value="{{ old('name') }}" required autofocus>

@if ($errors->has('name'))
<span role="alert">
<strong>{{ $errors->first('name') }}</strong>
</span>
@endif
</p>

<p>
<label for="email">{{ __('E-Mail Address') }}</label>

<input id="email" type="email" class="{{ $errors->has('email') ? ' is-invalid' : '' }}" name="email" value="{{ old('email') }}" required>

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
<label for="password-confirm">{{ __('Confirm Password') }}</label>

<input id="password-confirm" type="password" class="form-control" name="password_confirmation" required>
</p>

<p>
<button type="submit">
{{ __('Register') }}
</button>
</p>
</form>
@endsection
