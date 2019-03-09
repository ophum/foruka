<ul>
<a href="{{ route('containers.create') }}"><li>create container</li></a>
</ul>

<table border="1">
@foreach($containers as $c)
    <tr><td><a href="{{ route('containers.show', $c->id) }}">{{ $c->name }}</a></td><td>{{ $c->image->name }}</td></tr>
@endforeach
</table>