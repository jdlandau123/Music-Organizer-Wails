<script>
  import { onMount } from 'svelte';
  import { GetAlbums } from '../../wailsjs/go/main/App.js';
  import arrowIcon from '../assets/icons/arrow.svg';
  import { albumsStore } from '../stores/albumsStore.js';

  let allSelected = false;
  let currentSort = {column: 'Artist', direction: 'asc'};

  function getAlbums() {
    GetAlbums().then(a => $albumsStore = a);
  }

  function handleSelectAll(isAllSelected) {
    if (isAllSelected) {
      for (let a of $albumsStore) {
        a.IsOnDevice = true;
      }
      $albumsStore = [...$albumsStore];
    } else {
      getAlbums();
      currentSort = {column: 'Artist', direction: 'asc'};
    }
  }

  function updateSort(column) {
    if (column !== currentSort.column) {
      currentSort = {column, direction: 'asc'};
      return;
    }
    currentSort.column = column;
    if (currentSort.direction === 'asc') {
      currentSort.direction = 'desc';
    } else {
      currentSort.direction = 'asc';
    }
  }

  function sortTable(sortObj) {
    switch (sortObj.direction) {
      case 'asc':
        $albumsStore = $albumsStore.sort((a, b) => (a[sortObj.column] > b[sortObj.column]) ? 1 : ((b[sortObj.column] > a[sortObj.column]) ? -1 : 0));
        break;
      case 'desc':
        $albumsStore = $albumsStore.sort((a, b) => (a[sortObj.column] < b[sortObj.column]) ? 1 : ((b[sortObj.column] < a[sortObj.column]) ? -1 : 0));
        break;
      default:
        $albumsStore = $albumsStore.sort((a, b) => (a[sortObj.column] > b[sortObj.column]) ? 1 : ((b[sortObj.column] > a[sortObj.column]) ? -1 : 0));
        break;
    }
  }

  $: handleSelectAll(allSelected);

  $: sortTable(currentSort);

  onMount(() => getAlbums());
</script>

<table>
  <tr>
    <th>
      <input class="checkbox" type="checkbox" bind:checked={allSelected} />
    </th>
    <th on:click={() => updateSort('Artist')}>
      <div class="tblHeader">
        Artist 
        {#if currentSort.column === 'Artist'}
          <img class={currentSort.direction === 'desc' ? 'sortIcon rotated' : 'sortIcon'}
               src={arrowIcon} alt="arrow" />
        {/if}
      </div>
    </th>
    <th on:click={() => updateSort('Album')}>
      <div class="tblHeader">
        Album 
        {#if currentSort.column === 'Album'}
          <img class={currentSort.direction === 'desc' ? 'sortIcon rotated' : 'sortIcon'}
               src={arrowIcon} alt="arrow" />
        {/if}
      </div>
    </th>
    <th on:click={() => updateSort('FileFormat')}>
      <div class="tblHeader">
        File Format 
        {#if currentSort.column === 'FileFormat'}
          <img class={currentSort.direction === 'desc' ? 'sortIcon rotated' : 'sortIcon'}
               src={arrowIcon} alt="arrow" />
        {/if}
      </div>
    </th>
  </tr>
  {#each $albumsStore as row}
    <tr>
      <td>
        <input class="checkbox" type="checkbox" bind:checked={row.IsOnDevice} />
      </td>
      <td>{row.Artist}</td>
      <td>{row.Album}</td>
      <td>{row.FileFormat}</td>
    </tr>
  {/each}
</table>

<style>
  table {
    display: block;
    width: 60%;
    border-collapse: collapse;
    max-height: calc(100vh - 100px);
    overflow: auto;
  }

  tr {
    border-bottom: 1px solid lightgrey;
  }

  tr:last-child { 
    border-bottom: none; 
  }

  tr:hover {
    background-color: rgba(75, 163, 70, 0.3);
    cursor: pointer;
  }

  th, td {
    text-align: left;
    padding: 10px;
  }

  table td:nth-child(3) {
    width: 400px;
  }

  table td:nth-child(4) {
    width: 140px;
  }

  input.checkbox {
    transform: scale(1.3);
  }

  .tblHeader {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .sortIcon {
    transition-duration: 0.8s;
    transition-property: transform;
  }

  .rotated {
    transform: rotate(180deg);
  }
</style>
