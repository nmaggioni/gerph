function getStats(wantedBucket) {
    $('body').loadie();
    $.getJSON('/stats', function (stats) {
        $('body').loadie(0.25);

        $('#number-buckets').html(stats.bucketsNumber);
        $('#number-keys').html(stats.keysNumber);
        $('#disk-space-used').html((stats.diskUsedBytes / 1048576).toFixed(1));  // MB
        $('#disk-space-used').attr('title', stats.diskUsedBytes + ' bytes');
        $('#disk-space-free').html(stats.diskFreeBytes > 0 ? (stats.diskFreeBytes / 1073741824).toFixed(1) : 'N/A');  // GB
        $('#disk-space-free').attr('title', stats.diskFreeBytes > 0 ? stats.diskFreeBytes + ' bytes' : 'Not applicable on this OS');
        $('body').loadie(0.50);

        var bucket, key, i = 0, keyIndex;
        $('#buckets-sidebar li').remove();
        $('#keys-table-header .sub-header-text').text('Keys');
        $('#btn-add').off();  // Clear old handlers
        $('#btn-add').click(function() {
            var bucket = wantedBucket ? wantedBucket : prompt('In which bucket do you want to place the new key?');
            if (bucket) {
                var key = prompt("What should the new key be named?");
                if (key) {
                    var value = prompt("What should the value of the new key be?");
                    if (value) {
                        $.ajax({
                        type: "POST",
                        url: '/api/' + bucket + '/' + key,
                        data: { value: value },
                        success: function () {
                            alert('Key successfully added.');
                            getStats(bucket);
                        }
                    });
                    }
                }
            }
        });
        $('#keys-table tbody tr').remove();
        for (bucket in stats.keys) {
            $('#buckets-sidebar').append('<li title="' + bucket + '">' +
                '<a href="#">' + bucket +
                '<span class="btn-delete fa fa-trash" aria-hidden="true"></span>' +
                '</a>' +
                '</li>'
            );
            if (wantedBucket) {
                $('#keys-table-header .sub-header-text').text('Keys (' + wantedBucket + ')');
                if (wantedBucket !== bucket) {
                    continue;
                }
            }
            for (keyIndex in stats.keys[bucket]) {
                i++;
                key = stats.keys[bucket][keyIndex];
                $('#keys-table tbody').append('<tr>' +
                    '<td class="cell-index">' + i + '</td>' +
                    '<td class="cell-bucket">' + bucket + '</td>' +
                    '<td class="cell-key">' + key.key + '</td>' +
                    '<td class="cell-value">' + key.value + '</td>' +
                    '<td align="center"><i class="btn-edit fa fa-pencil" aria-hidden="true"></i></td>' +
                    '<td align="center"><i class="btn-delete fa fa-trash" aria-hidden="true"></i></td></tr>');
            }
        }
        $('body').loadie(0.75);

        $('.sidebar li').removeClass('active');
        if (wantedBucket) {
            $('.sidebar li').filter(function () {
                return $(this).text() === wantedBucket;
            }).addClass('active');
        } else {
            $('#overview-li').addClass('active');
        }

        $('#buckets-sidebar').off();  // Clear old handlers
        $('#buckets-sidebar').on('click', 'li a', function (e) {
            getStats($(this).text());
        });
        $('#buckets-sidebar').on('click', 'span.btn-delete', function (e) {
            e.stopPropagation();
            var bucket = $(this).closest('a').text();
            var yn = confirm('Do you really want to delete the bucket "' + bucket + '"?');
            if (yn) {
                $.ajax({
                    url: '/api/' + bucket,
                    type: 'DELETE',
                    success: function () {
                        alert('Bucket successfully deleted.');
                        getStats();
                    }
                });
            }
        });
        setTimeout(function () {
            $('body').loadie(1);
        }, 500);

        $('#keys-table .btn-edit').each(function () {
            var bucket = $(this).closest('tr').children('td.cell-bucket').text();
            var key = $(this).closest('tr').children('td.cell-key').text();
            var value = $(this).closest('tr').children('td.cell-value').text();
            $(this).click(function () {
                var newValue = prompt('What should the new value be for the key "' + key + '" from the bucket "' + bucket + '"?', value);
                if (newValue) {
                    $.ajax({
                        type: "POST",
                        url: '/api/' + bucket + '/' + key,
                        data: { value: newValue },
                        success: function () {
                            alert('Key successfully edited.');
                            getStats(wantedBucket);
                        }
                    });
                }
            });
        });
        $('#keys-table .btn-delete').each(function () {
            var bucket = $(this).closest('tr').children('td.cell-bucket').text();
            var key = $(this).closest('tr').children('td.cell-key').text();
            var value = $(this).closest('tr').children('td.cell-value').text();
            $(this).click(function () {
                var yn = confirm('Do you really want to delete the key "' + key + '" from the bucket "' + bucket + '"?');
                if (yn) {
                    $.ajax({
                        url: '/api/' + bucket + '/' + key,
                        type: 'DELETE',
                        success: function () {
                            alert('Key successfully deleted.');
                            getStats(wantedBucket);
                        }
                    });
                    console.log('In case you made a mistake... ' +
                        'BUCKET: "' + bucket + '", KEY: "' + key + '", VALUE: "' + value.replace(/"/g, '\\"') + '"');
                }
            });
        });

        if ($("#pagination-checkbox")[0].checked) paginate(wantedBucket);
    });
}

function paginate(wantedBucket) {
    if (paginator) {
        $('#paginator').remove();
    }
    paginator = $('#keys-table').simplePagination({
        perPage: wantedBucket ? 10 : 5,
        containerId: 'paginator',
        containerClass: '',
        previousButtonClass: 'btn btn-primary btn-paginator-prev',
        nextButtonClass: 'btn btn-primary btn-paginator-next',
        previousButtonText: '<',
        nextButtonText: '>',
        currentPage: 1
    });
}

var paginator;

(function () {
    getStats();

    $(document).keydown(function(e) {
        switch(e.which) {
            case 37: // left
                var prev = $('.btn-paginator-prev');
                if (prev) prev.trigger('click');
                break;

            case 38: // up
                var active = $('#buckets-sidebar li.active');
                if (active.length) {
                    if (active.prev().length) {
                        active.prev().find('a').trigger('click');
                    } else {
                        $('#overview-li').find('a').trigger('click');
                    }
                } else {
                    $('#buckets-sidebar li').last().find('a').trigger('click');
                }
                break;

            case 39: // right
                var next = $('.btn-paginator-next');
                if (next) next.trigger('click');
                break;

            case 40: // down
                var active = $('#buckets-sidebar li.active');
                if (active.length) {
                    if (active.next().length) {
                        active.next().find('a').trigger('click');
                    } else {
                        $('#overview-li').find('a').trigger('click');
                    }

                } else {
                    $('#buckets-sidebar li').first().find('a').trigger('click');
                }
                break;

            default:
                return;
        }
        e.preventDefault();
    });

    $("#pagination-checkbox").change(function() {
        if (this.checked) {
            paginate($('#buckets-sidebar li.active').length ? $('#buckets-sidebar li.active').find('a').text() : null);
        } else {
            $('#keys-table tbody tr').each(function(row) {
                $(this).show();
            });
            $('#paginator').remove();
        }
    });
})();