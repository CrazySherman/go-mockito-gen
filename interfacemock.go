import re
test_str =  'type Gladiator interface { \n \
    Spartacus(s TypeA) (TypeB, TypeC, error) \n \
    Critux() TypeD \n \
    Gannicus(TypeA) \n\
}\n  \
'

interface_name = ""
funcs = []
pattern_iname = 'type[\s]+([a-zA-Z]+)[\s]+interface'
pattern_func = '^[\s]*([a-zA-Z]+)(\([a-zA-Z\s,]*\))[\s]*(\(?[a-zA-Z\s,]+\)?)'
p1 = re.compile(pattern_iname)
p2 = re.compile(pattern_func)
for line in test_str.splitlines():
    print line
    m1 = p1.match(line)
    m2 = p2.match(line)
    if m1:
        interface_name = m1.group(1)
    elif m2:
        name, args, ret = m2.group(1,2,3)
#         print name, args, ret
        arg_list = []
        ret_list = []
        # parsing args
        tmp_list= args[1:-1].split(',')
        for tmp in tmp_list:
            tmp2 = tmp.strip().split()
            # non-arg function
            if not tmp2:
                continue
            if len(tmp2) == 1:
                arg_list.append(tmp2[0])
            elif len(tmp2) == 2: 
                arg_list.append(tmp2[1])     
            else:
                raise ValueError('parsing failed for argument: %s' % tmp)
                
        # if return type exist, parse rets
        if ret:    
            tmp_list = ret.strip('(').strip(')').split(',')
            for tmp in tmp_list:
                tmp2 = tmp.strip().split()
                if not tmp2:
                    continue
                elif len(tmp2) == 1: 
                    ret_list.append(tmp2[0])
                elif len(tmp2) == 2:
                    ret_list.append(tmp2[1])
                else:
                    raise ValueError('parsing failed for return type: %s' % tmp)
                    
        func_descriptor = {
            'name': name,
            'args': arg_list,
            'rets': ret_list
        }
        funcs.append(func_descriptor)
    else:
        continue
# print 'Interface name: ', interface_name
# print 'function descriptors: ', funcs

print '################### GENERATED MOCK CLASS ########################'
print 'type ' + interface_name + 'Mock struct {'
print '    mock.Mock'
print '}'

for func in funcs:
    print ''
    argstr = ''
    arg_names = ''
    retstr = ''
    ret_names = ''
    # generating args
    for i, arg_type in enumerate(func['args']):
        argstr = argstr + 'arg' + str(i) + ' ' + arg_type + ', '
        arg_names = arg_names + 'arg' + str(i) + ', '
    # if function has any args
    if argstr:
        argstr = argstr[:-2]
        arg_names = arg_names[:-2]
        
    # generating rets
    for i, ret_type in enumerate(func['rets']):
        retstr = retstr + ret_type + ', '
        if ret_type == 'error':
            ret_names = ret_names + 'args.Error(' + str(i) + ').(' + ret_type + '), '
        else:
            ret_names = ret_names + 'args.Get(' + str(i) + ').(' + ret_type + '), '
        
    # if there is any rets
    if retstr:
        retstr = retstr[:-2]
        ret_names = ret_names[:-2]
        
    if len(func['rets']) > 1:
        retstr = '(' + retstr + ')'
    
    print 'func(m ', interface_name, 'Mock)' + func['name'] + '(' + argstr + ') ' + retstr + ' {'
    print '    args := m.Called(' + arg_names + ')'
    # if there is anything to return
    if retstr:
        print '    return ' + ret_names
    print '}'
        
        