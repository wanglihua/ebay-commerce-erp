/**
 * Sometimes for quick navigation, it can be useful to allow an end user to
 * enter which page they wish to jump to manually. This paging control uses a
 * text input box to accept new paging numbers (arrow keys are also allowed
 * for), and four standard navigation buttons are also presented to the end
 * user.
 *
 *  @name Navigation with text input
 *  @summary Shows an input element into which the user can type a page number
 *  @author [Allan Jardine](http://sprymedia.co.uk)
 *  @author [Gordey Doronin](http://github.com/GDoronin)
 *
 *  @example
 *    $(document).ready(function() {
 *        $('#example').dataTable( {
 *            "pagingType": "input"
 *        } );
 *    } );
 */

(function ($) {
    function calcDisableClasses(oSettings) {
        var start = oSettings._iDisplayStart;
        var length = oSettings._iDisplayLength;
        var visibleRecords = oSettings.fnRecordsDisplay();
        var all = length === -1;

        // Gordey Doronin: Re-used this code from main jQuery.dataTables source code. To be consistent.
        var page = all ? 0 : Math.ceil(start / length);
        var pages = all ? 1 : Math.ceil(visibleRecords / length);

        var disableFirstPrevClass = (page > 0 ? '' : oSettings.oClasses.sPageButtonDisabled);
        var disableNextLastClass = (page < pages - 1 ? '' : oSettings.oClasses.sPageButtonDisabled);

        return {
            'first': disableFirstPrevClass,
            'previous': disableFirstPrevClass,
            'next': disableNextLastClass,
            'last': disableNextLastClass
        };
    }

    function calcCurrentPage(oSettings) {
        return Math.ceil(oSettings._iDisplayStart / oSettings._iDisplayLength) + 1;
    }

    function calcPages(oSettings) {
        return Math.ceil(oSettings.fnRecordsDisplay() / oSettings._iDisplayLength);
    }

    var firstClassName = 'first';
    var previousClassName = 'previous';
    var nextClassName = 'next';
    var lastClassName = 'last';

    var paginateClassName = 'paginate';
    var paginateOfClassName = 'paginate_of';
    var paginatePageClassName = 'paginate_page';
    var paginateInputClassName = 'paginate_input';

    $.fn.dataTableExt.oPagination.input = {
        'fnInit': function (oSettings, nPaging, fnCallbackDraw) {
            var nUl = document.createElement('ul');
            nUl.className = 'pagination';
            nPaging.appendChild(nUl);

            var nFirst = document.createElement('li');
            var nPrevious = document.createElement('li');
            var nNext = document.createElement('li');
            var nLast = document.createElement('li');
            var nInput = document.createElement('input');
            var nPage = document.createElement('span');
            var nOf = document.createElement('span');

            var language = oSettings.oLanguage.oPaginate;
            var classes = oSettings.oClasses;

            var nFirstA = document.createElement('a');
            //nFirstA.setAttribute("href", "#");
            nFirstA.innerHTML = language.sFirst;
            nFirst.appendChild(nFirstA);

            var nPreviousA = document.createElement('a');
            //nPreviousA.setAttribute("href", "#");
            nPreviousA.innerHTML = language.sPrevious;
            nPrevious.appendChild(nPreviousA);

            var nNextA = document.createElement('a');
            //nNextA.setAttribute("href", "#");
            nNextA.innerHTML = language.sNext;
            nNext.appendChild(nNextA);

            var nLastA = document.createElement('a');
            //nLastA.setAttribute("href", "#");
            nLastA.innerHTML = language.sLast;
            nLast.appendChild(nLastA);

            nFirst.className = firstClassName + ' ' + classes.sPageButton;
            nPrevious.className = previousClassName + ' ' + classes.sPageButton;
            nNext.className = nextClassName + ' ' + classes.sPageButton;
            nLast.className = lastClassName + ' ' + classes.sPageButton;

            nOf.className = paginateOfClassName;
            nPage.className = paginatePageClassName;
            nInput.className = paginateInputClassName;
            nInput.setAttribute("style", "width:50px;text-align:center;");

            if (oSettings.sTableId !== '') {
                nPaging.setAttribute('id', oSettings.sTableId + '_' + paginateClassName);
                nFirst.setAttribute('id', oSettings.sTableId + '_' + firstClassName);
                nPrevious.setAttribute('id', oSettings.sTableId + '_' + previousClassName);
                nNext.setAttribute('id', oSettings.sTableId + '_' + nextClassName);
                nLast.setAttribute('id', oSettings.sTableId + '_' + lastClassName);
            }

            nInput.type = 'text';
            nPage.innerHTML = '&nbsp;&nbsp;&nbsp;第 ';

            nUl.appendChild(nFirst);
            nUl.appendChild(nPrevious);
            nUl.appendChild(nNext);
            nUl.appendChild(nLast);

            nUl.appendChild(nPage);
            nUl.appendChild(nInput);
            nUl.appendChild(nOf);

            $(nFirst).click(function () {
                var iCurrentPage = calcCurrentPage(oSettings);
                if (iCurrentPage !== 1) {
                    oSettings.oApi._fnPageChange(oSettings, 'first');
                    fnCallbackDraw(oSettings);
                }
            });

            $(nPrevious).click(function () {
                var iCurrentPage = calcCurrentPage(oSettings);
                if (iCurrentPage !== 1) {
                    oSettings.oApi._fnPageChange(oSettings, 'previous');
                    fnCallbackDraw(oSettings);
                }
            });

            $(nNext).click(function () {
                var iCurrentPage = calcCurrentPage(oSettings);
                if (iCurrentPage !== calcPages(oSettings)) {
                    oSettings.oApi._fnPageChange(oSettings, 'next');
                    fnCallbackDraw(oSettings);
                }
            });

            $(nLast).click(function () {
                var iCurrentPage = calcCurrentPage(oSettings);
                if (iCurrentPage !== calcPages(oSettings)) {
                    oSettings.oApi._fnPageChange(oSettings, 'last');
                    fnCallbackDraw(oSettings);
                }
            });

            $(nInput).keyup(function (e) {
                // 38 = up arrow, 39 = right arrow
                if (e.which === 38 || e.which === 39) {
                    this.value++;
                }
                // 37 = left arrow, 40 = down arrow
                else if ((e.which === 37 || e.which === 40) && this.value > 1) {
                    this.value--;
                }

                if (this.value === '' || this.value.match(/[^0-9]/)) {
                    /* Nothing entered or non-numeric character */
                    this.value = this.value.replace(/[^\d]/g, ''); // don't even allow anything but digits
                    return;
                }

                var iNewStart = oSettings._iDisplayLength * (this.value - 1);
                if (iNewStart < 0) {
                    iNewStart = 0;
                }
                if (iNewStart >= oSettings.fnRecordsDisplay()) {
                    iNewStart = (Math.ceil((oSettings.fnRecordsDisplay()) / oSettings._iDisplayLength) - 1) * oSettings._iDisplayLength;
                }

                oSettings._iDisplayStart = iNewStart;
                fnCallbackDraw(oSettings);
            });

            // Take the brutal approach to cancelling text selection.
            $('li', nPaging).bind('mousedown', function () {
                return false;
            });
            $('li', nPaging).bind('selectstart', function () {
                return false;
            });

            // If we can't page anyway, might as well not show it.
            var iPages = calcPages(oSettings);
            if (iPages <= 1) {
                $(nPaging).hide();
            }
        },

        'fnUpdate': function (oSettings) {
            if (!oSettings.aanFeatures.p) {
                return;
            }

            var iPages = calcPages(oSettings);
            var iCurrentPage = calcCurrentPage(oSettings);

            var an = oSettings.aanFeatures.p;

            if (iPages <= 1) // hide paging when we can't page
            {
                $(an).hide();
                return;
            }

            var disableClasses = calcDisableClasses(oSettings);

            $(an).show();

            // Enable/Disable `first` button.
            var ul_pager = $(an).children('ul');
            ul_pager.children('.' + firstClassName)
                .removeClass(oSettings.oClasses.sPageButtonDisabled)
                .addClass(disableClasses[firstClassName]);

            // Enable/Disable `prev` button.
            ul_pager.children('.' + previousClassName)
                .removeClass(oSettings.oClasses.sPageButtonDisabled)
                .addClass(disableClasses[previousClassName]);

            // Enable/Disable `next` button.
            ul_pager.children('.' + nextClassName)
                .removeClass(oSettings.oClasses.sPageButtonDisabled)
                .addClass(disableClasses[nextClassName]);

            // Enable/Disable `last` button.
            ul_pager.children('.' + lastClassName)
                .removeClass(oSettings.oClasses.sPageButtonDisabled)
                .addClass(disableClasses[lastClassName]);

            // Paginate of N pages text
            ul_pager.children('.' + paginateOfClassName).html(' / ' + iPages + ' 页 ');

            // Current page numer input value
            ul_pager.children('.' + paginateInputClassName).val(iCurrentPage);
        }
    };
})(jQuery);


function dtFormatDate(data) {
    if (!data) {
        return "";
    }
    data = data.replace(/T|Z/g, ' ');
    return stringToDate(data).format("yyyy-MM-dd");
}

function dtFormatDateTime(data) {
    if (!data) {
        return "";
    }
    data = data.replace(/T|Z/g, ' ');
    return stringToDate(data).format("yyyy-MM-dd hh:mm");
}

function dtFormatDateTime2(data) {
    if (!data) {
        return "";
    }
    data = data.replace(/T|Z/g, ' ');
    return stringToDate(data).format("yyyy-MM-dd hh:mm:ss");
}

function dtFormatFloat1(data) {
    if ((!data) && (data !== 0)) {
        return "";
    }
    return $.formatNumber(data, {format: "###0.0", locale: "cn"});
}

function dtFormatFloat(data) {
    if ((!data) && (data !== 0)) {
        return "";
    }
    return $.formatNumber(data, {format: "###0.00", locale: "cn"});
}

function dtFormatFloat3(data) {
    if ((!data) && (data !== 0)) {
        return "";
    }
    return $.formatNumber(data, {format: "###0.000", locale: "cn"});
}

function dtFormatFloat4(data) {
    if ((!data) && (data !== 0)) {
        return "";
    }
    return $.formatNumber(data, {format: "###0.0000", locale: "cn"});
}

function dtFormatInt(data) {
    if ((!data) && (data !== 0)) {
        return "";
    }
    return $.formatNumber(data, {format: "###0", locale: "cn"});
}

$.messager.model = {
    ok: {text: '<i class="ace-icon fa fa-check"></i>确认', classed: 'btn-success btn-sm'},
    cancel: {text: '<i class="ace-icon fa fa-times"></i>取消', classed: 'btn-default btn-sm'}
};

$.msg = {
    error: function (text) {
        $.messager.alert("错误", '<div class="huaduConfirm"><i class="fa fa-times-circle fa-5x red"></i></br></br>' + text + '</div>');
    },
    //waring: function (text) {
    //    $.messager.alert("警告", '<div class="huaduConfirm"><i class="fa fa-exclamation-triangle fa-5x orange2"></i></br></br>' + text + '</div>');
    //},
    alert: function (text) {
        $.messager.alert("提示", '<div class="huaduConfirm"><i class="fa fa-exclamation-triangle fa-5x orange2"></i></br></br>' + text + '</div>');
    },
    confirm: function (text, fn) {
        $.messager.confirm("确认", '<div class="huaduConfirm"><i class="fa fa-question-circle fa-5x orange2"></i></br></br> ' + text + '</div>', fn);
    },
    tip: function (text) {
        $.messager.alert("确认", '<div class="huaduConfirm"><i class="fa fa-exclamation-circle fa-5x blue"></i></br></br>' + text + '</div>');
    },
    show: function (text) {
        $.messager.popup("<span style='color:green;font-weight:bold;'>" + text + "</span>");
    }
};

window.onresize = function () {
    $("#huaduCover").height($(window).height()).css("line-height", $(window).height() + "px").width($(window).width());
};

/*
 $(document).ajaxStart(function () {
 showMask("数据加载中，请稍候");
 }).ajaxStop(function () {
 hideMask();
 });
 */

//显示遮照层
function showMask(text) {
    //$('body').css("overflow", "hidden")
    $("#huaduCover").show();
    $("#huaduProgressBar").html(text);
}

//隐藏遮照层
function hideMask() {
    //$('body').css("overflow", "")
    $("#huaduCover").hide();
}

function handleJQueryAjaxError(xhr, textStatus, error) {
    hideMask();
    if (textStatus === 'timeout') {
        $.msg.error('服务器超时！');
    }
    else {
        $.msg.error('访问服务器出错，请重试！' + error);
    }
}

$.ajaxSetup(
    {
        type: "POST",
        "timeout": 600000,
        "error": handleJQueryAjaxError
    }
);

//扩展dataTable默认值
$.extend($.fn.dataTable.defaults, {
    "sDom": "<'row'<'col-xs-6'l><'col-xs-6'f>r><'row'<'col-xs-12't>><'row'<'col-xs-6'i><'col-xs-6'p>>",
    //"bAutoWidth": true,
    //"sScrollX": "100%",
    //"bScrollCollapse": true,
    "bAutoWidth": false,
    "bFilter": true,
    "bSort": true,
    "bInfo": true,
    "bLengthChange": true,
    "bServerSide": true,
    "sServerMethod": "POST",
    "bProcessing": false,
    "bPaginate": true,
    "pagingType": "input",
    "lengthMenu": [[5, 10, 15, 25, 50, 100, 300, 500, 1000, 2000, 3000, 5000, 10000], [5, 10, 15, 25, 50, 100, 300, 500, 1000, 2000, 3000, 5000, 10000]],
    "iDisplayLength": 10,
    "sPaginationType": "full_numbers",
    "oLanguage": {
        "sLengthMenu": "每页显示 _MENU_ 条记录",
        "sZeroRecords": "没有匹配数据",
        "sInfo": "当前显示第 _START_ 至 _END_ 条，共 _TOTAL_ 条记录",
        "sInfoEmpty": "显示第 0 至 0 项结果，共 0 项",
        "sInfoFiltered": "(从 _MAX_ 条记录中过滤)",
        "sInfoPostFix": "",
        "sSearch": "搜索:",
        "sUrl": "",
        "sEmptyTable": "没有查询到任何数据！",
        "sLoadingRecords": "载入中...",
        //"sInfoThousands": ',',
        "sThousands": '',
        "oPaginate": {
            "sFirst": "首页",
            "sPrevious": "上页",
            "sNext": "下页",
            "sLast": "末页"
        },
        "oAria": {
            "sSortAscending": ": 以升序排列此列",
            "sSortDescending": ": 以降序排列此列"
        }
    },
    "fnServerData": function (sSource, aoData, fnCallback) {
        showMask("数据加载中，请稍候");
        $.ajax(
            {
                "dataType": 'json',
                "url": sSource,
                "data": aoData,
                "success": function (data) {
                    hideMask();
                    if (data.success === false) {
                        $.msg.error(data.message);
                        return;
                    }
                    fnCallback(data);
                }
            });
    },
    "initComplete": function (settings) {
        var dataTable = this;
        var searchHTML = '<label class="dataTable-custom-search"><span>搜索:</span><input type="search" class="form-control input-sm" placeholder="输入搜索内容按回车" aria-controls="' + dataTable.selector.substring(1) + '"></label>';
        //快捷操作的HTML DOM
        $(dataTable.selector + '_wrapper .dataTables_filter').append(searchHTML);

        //重写搜索事件
        $(dataTable.selector + '_wrapper .dataTables_filter input').bind('keyup',
            function (e) {
                if (e.keyCode === 13 || (e.keyCode === 8 && (this.value.length === 0))) {
                    dataTable.api().search(this.value).draw();
                }
            });
    }
});

function initDataTable(table, options) {

    var htmlTable = $(table);
    var dataTable = htmlTable.dataTable(options);

    function checkAll() {
        var thCheckBox = this;
        $(table).find('tbody tr>td:first-child input:checkbox').each(function () {
            var tdCheckBox = this;
            tdCheckBox.checked = thCheckBox.checked;
            var tableRow = $(this).closest('tr');
            if (tdCheckBox.checked) {
                tableRow.addClass('selected');
            } else {
                tableRow.removeClass('selected');
            }
        });
    }

    htmlTable.on('click', 'thead th input.checkall', checkAll);

    $(table + '_wrapper').on('click', 'div.dataTables_scroll div.dataTables_scrollHead thead th input.checkall', checkAll);

    /*
     htmlTable.on('click', 'tbody tr', function () {
     var tableRow = $(this);
     var checkbox = tableRow.find('td:first-child input:checkbox')[0];
     if (checkbox != undefined) {
     checkbox.checked = !checkbox.checked;
     if (checkbox.checked) {
     tableRow.addClass('selected');
     } else {
     tableRow.removeClass('selected');
     }
     }
     });
     */

    htmlTable.on('click', 'tbody td:first-child input:checkbox', function (event) {
        event.stopPropagation();
        var checkbox = $(this)[0];
        var tableRow = $(this).closest('tr');
        if (checkbox.checked) {
            tableRow.addClass('selected');
        } else {
            tableRow.removeClass('selected');
        }

        tableRow.removeClass('hilight');
    });

    return dataTable;
}

function initSingleSelectDataTable(table, options) {
    var htmlTable = $(table);
    htmlTable.dataTable(options);

    htmlTable.on('click', 'tbody tr', function () {
        var tableRow = $(this);
        var checkbox = tableRow.find('td:first-child input:checkbox')[0];
        if (checkbox !== undefined) {
            clearSelected(htmlTable);
            checkbox.checked = true;
            if (checkbox.checked) {
                tableRow.addClass('selected');
            }
        }
    });

    htmlTable.on('click', 'tbody td:first-child input:checkbox', function (event) {
        event.stopPropagation();
        var checkbox = $(this)[0];
        clearSelected(htmlTable, checkbox);
        var tableRow = $(this).closest('tr');
        if (checkbox.checked) {
            tableRow.addClass('selected');
        }
    });

    function clearSelected(htmlTable, exceptedCheckbox) {
        htmlTable.find('tbody tr').removeClass('selected');
        htmlTable.find('tbody td:first-child input:checkbox').each(function () {
            var checkbox = $(this)[0];
            if (exceptedCheckbox !== checkbox) {
                checkbox.checked = false;
            }
        });
    }
}

//重新加载数据，但是分页页面还是停留在原来的分页号
function reloadList(list) {
    var tables = $(list).dataTable().api();
    tables.draw(false);
}

//重新加载数据，分页页面变换为 第一页
function refreshList(list) {
    var tables = $(list).dataTable().api();
    tables.draw(true);
}

function trim(str) { //删除左右两端的空格
    return str.replace(/(^\s*)|(\s*$)/g, "");
}

//bootstrap modal options
var modalOptions = {backdrop: 'static', keyboard: false};

$.extend($.validator.messages, {
    required: "必填",
    remote: "请修正该字段",
    email: "请输入正确格式的电子邮件",
    url: "请输入合法的网址",
    date: "请输入合法的日期",
    dateISO: "请输入合法的日期(ISO).",
    number: "请输入合法的数字",
    digits: "只能输入整数",
    creditcard: "请输入合法的信用卡号",
    equalTo: "请再次输入相同的值",
    accept: "请输入拥有合法后缀名的字符串",
    maxlength: jQuery.validator.format("请输入一个长度最多是 {0} 的字符串"),
    minlength: jQuery.validator.format("请输入一个长度最少是 {0} 的字符串"),
    rangelength: jQuery.validator.format("请输入一个长度介于 {0} 和 {1} 之间的字符串"),
    range: jQuery.validator.format("请输入一个介于 {0} 和 {1} 之间的值"),
    max: jQuery.validator.format("请输入一个最大为 {0} 的值"),
    min: jQuery.validator.format("请输入一个最小为 {0} 的值")
});

function validate(form, rules, errorPlacement) {
    if (errorPlacement === undefined) {
        form.validate({
            submitHandler: function () {
            },
            rules: rules,
            highlight: function (e) {
                $(e).closest('.form-group').addClass('has-error');
            },
            success: function (e) {
                $(e).closest('.form-group').removeClass('has-error');
                $(e).remove();
            },
            errorElement: 'div',
            errorClass: 'help-block inline err-div',
            focusInvalid: false,
            ignore: "",
            errorPlacement: function (error, element) {
                //var parent = element.parent("div");
                var parent = element.parent("div[class!='input-group']");
                //if(!parent){
                //    parent = element.parent().parent();
                //}
                error.appendTo(parent);
            }
        });
    }
    else {
        form.validate({
            submitHandler: function () {
            },
            rules: rules,
            highlight: function (e) {
                $(e).closest('.form-group').addClass('has-error');
            },
            success: function (e) {
                $(e).closest('.form-group').removeClass('has-error');
                $(e).remove();
            },
            errorElement: 'div',
            errorClass: 'help-block inline err-div',
            focusInvalid: false,
            ignore: "",
            errorPlacement: errorPlacement
        });
    }

    return form.valid();
}

$.fn.clearValidation = function () {
    $(this).find("div.has-error").removeClass('has-error');
    $(this).find("div.err-div").remove();
};

/*
 function getCookie(name) {
 var cookieValue = null;
 if (document.cookie && document.cookie != '') {
 var cookies = document.cookie.split(';');
 for (var i = 0; i < cookies.length; i++) {
 var cookie = jQuery.trim(cookies[i]);
 // Does this cookie string begin with the name we want?
 if (cookie.substring(0, name.length + 1) == (name + '=')) {
 cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
 break;
 }
 }
 }
 return cookieValue;
 }

 var csrftoken = getCookie('csrftoken');

 var csrftoken = Cookies.get('csrftoken');
 */


/*
 function csrfSafeMethod(method) {
 // these HTTP methods do not require CSRF protection
 return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
 }

 $.ajaxSetup({
 beforeSend: function (xhr, settings) {
 if (!csrfSafeMethod(settings.type) && !this.crossDomain) {
 var csrftoken = Cookies.get('csrftoken');
 xhr.setRequestHeader("X-CSRFToken", csrftoken);
 }
 }
 });

 */

function isEmptyStr(str) {
    if (!str || str.length === 0 || str === "" || typeof str === 'undefined' || !/[^\s]/.test(str) || /^\s*$/.test(str) || str.replace(/\s/g, "") === "") {
        return true;
    } else {
        return false;
    }
}

