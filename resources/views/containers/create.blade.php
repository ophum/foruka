
<form action="{{ route('containers.store') }}" method="post">
{{csrf_field()}}
    <input type="text" id="container_name" name="container_name" value="{{ old('container_name') }}">
    <select id="container_image" name="container_image">
        <option value="ubuntu">Ubuntu</option>
        <option value="debian">Debian</option>
    </select>
    <input type="submit" value="make">
</form>