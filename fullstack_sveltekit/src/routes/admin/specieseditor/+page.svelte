<script lang="ts">
    import {NativeSelect} from '@svelteuidev/core';
    import {onMount} from "svelte";

    let pictures = [], species = [];
    let selectedSpecie, selectedSpecieName, picturesInSpecie;
    let newSpecie = {
        name: '',
        scientific_name: '',
        threat: '',
        info_page: '',
        description: '',
    }

    function handleSpecieSelectChange(evt) {
        if (pictures.length > 0 && species.length > 0) {
            selectedSpecie = species.filter(s => s.name === selectedSpecieName)[0];
            picturesInSpecie = selectedSpecie.pictures.map(p => p.id);
            console.log(selectedSpecie.name);
        }
    }

    $: selectedSpecieName, handleSpecieSelectChange()

    async function toggleInSpecie(picture) {
        console.log(picture, selectedSpecie)
        const action = selectedSpecie.pictures.map(p => p.id).includes(picture.id) ? 'remove' : 'add';

        fetch('/api/species', {
            method: 'PUT',
            body: JSON.stringify({
                pictureId: picture.id,
                specieId: selectedSpecie.id,
                action: action
            })
        }).then(res => {
            if (res.status === 200) {
                loadSpecies(selectedSpecieName);
            }
        })
    }

    function loadSpecies(updated) {
        fetch('/api/species').then(res => {
            res.json().then(data => {
                $: species = data.species;
                if (species.length > 0) {
                    if (updated) {
                        selectedSpecie = species.filter(a => a.name === updated)[0];
                    } else {
                        selectedSpecie = species[0];
                    }
                    selectedSpecieName = selectedSpecie.name;
                    picturesInSpecie = selectedSpecie.pictures.map(p => p.id);
                }
            });
        });
    }

    onMount(() => {
        loadSpecies(undefined);

        fetch('/api/pictures').then(res => {
            res.json().then(data => {
                $: pictures = data.pictures;
            });
        });
    })

</script>

<main>
    <h2>Nouvelle espèce</h2>
    {#if pictures.length > 0 && species.length > 0}
        <NativeSelect data={species.map(a => a.name)}
                      bind:value={selectedSpecieName}
                      label="Espece a éditer"/>
        <div id="picture-display">
            {#each pictures as picture}
                <img src={'/api/pictures/' + picture.path + '?type=thumb'}
                     on:click={() => toggleInSpecie(picture)}
                     class:selected={picturesInSpecie?.includes(picture.id)}>
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
