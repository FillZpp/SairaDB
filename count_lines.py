#! /usr/bin/env python

import os
import sys
import crash_on_ipy

fext_list = []

class CountLines(object):
    lines_all = 0
    lines_not_null = 0
    
    def count(self, d):
        #raise KeyError(1)
        if os.path.isdir(d):
            for i in os.listdir(d):
                self.count(d + '/' + i)
        else:
            for i in fext_list:
                if d.endswith(i):
                    with open(d, 'r') as fi:
                        li = fi.readlines()
                    for i in li:
                        i = i.strip()
                        if i:
                            self.lines_not_null += 1
                        self.lines_all += 1
                    break
                
def main ():
    argc = len(sys.argv)
    if argc == 1:
        print('Please appoint')
        return
    for i in range(1, argc):
        fext_list.append(sys.argv[i])

    cl = CountLines()
    cl.count(os.getcwd())
    print('All lines num: %d' % cl.lines_all)
    print('Not null lines num: %d' % cl.lines_not_null)


if __name__ == '__main__':
    main()

