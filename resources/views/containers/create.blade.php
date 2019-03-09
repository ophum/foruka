
<form action="{{ route('containers.store') }}" method="post">
{{csrf_field()}}
    <label for="container_name">
        name
    </label>
    <input type="text" id="container_name" name="container_name" value="{{ old('container_name') }}">
    <label for="container_image">
        image
    </label>
    <select id="container_image" name="container_image">
        <option value="ubuntu">Ubuntu</option>
        <option value="debian">Debian</option>
    </select>
    <input type="submit" value="create">
</form>