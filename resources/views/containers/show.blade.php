@extends('layouts.app')

@section('content')
<h3>{{ $container->name }}</h3>
<hr>
<p>this container image is {{ $container->image->name }}</p>
@endsection