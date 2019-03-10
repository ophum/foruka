@extends('layouts.app')

@section('content')
<form action="{{ route('containers.store') }}" method="post">
{{csrf_field()}}
    <label for="name">
        name
    </label>
    <input type="text" id="name" name="name" value="{{ old('name') }}">
    <label for="image">
        image
    </label>
    <select id="image" name="image">
    @foreach($images as $image)
        <option value="{{ $image->name }}">{{ $image->name }}</option>
    @endforeach
    </select>
    <input type="submit" value="create">
</form>
@endsection