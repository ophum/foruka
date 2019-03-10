@extends('layouts.app')

@section('content')

<table class="table table-striped">
    <thead>
        <tr>
            <th>Container Name</th>
            <th>Image Name</th>
            <th>Status</th>
        </tr>
    </thead>

    <tbody>
@foreach($containers as $c)
        <tr>
            <td>
                <a href="{{ route('containers.show', $c->id) }}"><div>{{ $c->name }}</div></a>
            </td>
            <td>
                {{ $c->image->name }}
            </td>
            <td>

            </td>
        </tr>
@endforeach
    </tbody>
</table>

@endsection