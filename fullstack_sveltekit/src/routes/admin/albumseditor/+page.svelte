<script lang="ts">
    import {Button, NativeSelect, TextInput} from '@svelteuidev/core';
    import {onMount} from "svelte";
    import {getHeaders} from "$lib/services/adminHeaders";

    let pictures = [], albums = [];

    let selectedAlbum,// = albums[0],
        selectedAlbumName,// = selectedAlbum.name,
        picturesInAlbum;// = selectedAlbum.pictures.map(p => p.id);

    let newAlbumName, newAlbumDesc;

    function handleAlbumSelectChange(evt) {
        if (pictures.length > 0 && albums.length > 0) {
            selectedAlbum = albums.filter(a => a.name === selectedAlbumName)[0];
            picturesInAlbum = selectedAlbum.pictures.map(p => p.id);
            console.log(selectedAlbum.name);
        }
    }

    $: selectedAlbumName, handleAlbumSelectChange()

    async function toggleInAlbum(picture) {
        console.log(picture, selectedAlbum)
        const action = selectedAlbum.pictures.map(p => p.id).includes(picture.id) ? 'remove' : 'add';

        fetch('/api/albums', {
            method: 'PUT',
            headers: getHeaders(),
            body: JSON.stringify({
                pictureId: picture.id,
                albumId: selectedAlbum.id,
                action: action
            })
        }).then(res => {
            if (res.status === 200) {
                loadAlbums(selectedAlbumName);
            }
        })
    }

    function loadAlbums(updated) {
        fetch('/api/albums').then(res => {
            res.json().then(data => {
                $: albums = data.albums;
                if (albums.length > 0) {
                    if (updated) {
                        selectedAlbum = albums.filter(a => a.name === updated)[0];
                    } else {
                        selectedAlbum = albums[0];
                    }
                    console.log('selectedalbum', selectedAlbum);
                    selectedAlbumName = selectedAlbum.name;
                    picturesInAlbum = selectedAlbum.pictures.map(p => p.id);
                }
            });
        });
    }

    function newAlbum() {
        console.log(newAlbumName, newAlbumDesc);

        fetch('/api/albums', {
            method: 'POST',
            headers: getHeaders(),
            body: JSON.stringify({
                name: newAlbumName,
                description: newAlbumDesc
            })
        }).then(res => {
            console.log(res);
            loadAlbums(undefined);
        })

    }

    onMount(() => {
        loadAlbums(undefined);

        fetch('/api/pictures').then(res => {
            res.json().then(data => {
                $: pictures = data.pictures;
            });
        });
    })

</script>

<main>
    <TextInput placeholder="Nom" label="Nom" bind:value={newAlbumName}/>
    <TextInput placeholder="Description" label="Description" bind:value={newAlbumDesc}/>
    <Button on:click={newAlbum}>Nouvel album</Button>

    {#if pictures.length > 0 && albums.length > 0}
        <NativeSelect data={albums.map(a => a.name).sort()}
                      bind:value={selectedAlbumName}
                      label="Album a éditer"/>
        <div id="picture-display">
            {#each pictures as picture}
                <img src={'/api/pictures/' + picture.path + '?type=thumb'}
                     on:click={() => toggleInAlbum(picture)}
                     class:selected={picturesInAlbum?.includes(picture.id)}>
            {/each}
        </div>
    {/if}
</main>

<style lang="scss">
  main {
    position: absolute;
    top: 7vh;
    left: 10vw;
    width: 80vw;

    #picture-display {
      width: 100%;
      margin-top: 20px;
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;

      img {
        height: 150px;
        margin: 5px;
        border: 5px solid #ddd;
      }
    }
  }

  :global(.selected) {
    border-color: red !important;
  }
</style>
