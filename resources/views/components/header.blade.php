<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
  <a class="navbar-brand" href="{{ route('home') }}">
    foruka
  </a>
  <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbar" aria-controls="navbar" aria-expanded="false" aria-label="Toggle navigation">
    <span class="navbar-toggler-icon"></span>
  </button>
  <div class="collapse navbar-collapse" id="navbar">
    <ul class="navbar-nav">

        <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="navbarContainer" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
            Container
            </a>
            <div class="dropdown-menu" aria-labelledby="navbarContainer">
            <a class="dropdown-item" href="{{ route('containers.index') }}">Index</a>
            <a class="dropdown-item" href="{{ route('containers.create') }}">Launch</a>
            </div>
        </li>
    </ul>
  </div>
</nav>