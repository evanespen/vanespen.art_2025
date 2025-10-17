<div id="side-title">EVRARD VAN ESPEN</div>
<div id="side-subtitle">Photographie</div>

<style lang="scss">
  @import "$src/color.scss";
  @import "$src/fonts.scss";

  #side-title {
    @include f-h;
    color: $subcolor;
    opacity: .1;
    z-index: 1000000;
    font-size: 6em;

    display: flex;
    justify-content: center;

    position: fixed;
    width: 100vh;
    top: 95vh;
    left: 4.5vh;

    transform: rotate(-90deg);
    transform-origin: left center;

    pointer-events: none;

    @media (max-width: 800px) {
      font-size: 4em;
    }
  }

  #side-subtitle {
    @include f-h;
    font-size: 4em;
    color: $subcolor;
    opacity: .1;
    z-index: 1000000;

    position: fixed;
    width: 100vh;
    top: 96vh;
    right: 3vh;

    display: flex;
    justify-content: center;

    transform: rotate(90deg);
    transform-origin: right center;

    pointer-events: none;

    @media (max-width: 800px) {
      font-size: 3em;
    }
  }
</style>
