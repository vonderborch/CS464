
/*
  WARNING: THIS FILE IS AUTO-GENERATED. DO NOT MODIFY.

  This file was generated from .idl using "rtiddsgen".
  The rtiddsgen tool is part of the RTI Connext distribution.
  For more information, type 'rtiddsgen -help' at a command shell
  or consult the RTI Connext manual.
*/
    
import com.rti.dds.typecode.*;


public class TopFunctionTypeCode {
    public static final TypeCode VALUE = getTypeCode();

    private static TypeCode getTypeCode() {
        TypeCode tc = null;
        int i=0;
        StructMember sm[] = new StructMember[6];

        sm[i]=new StructMember("username",false,(short)-1,false,(TypeCode)new TypeCode(TCKind.TK_STRING,255),0,false); i++;
        sm[i]=new StructMember("hostname",false,(short)-1,false,(TypeCode)new TypeCode(TCKind.TK_STRING,255),1,false); i++;
        sm[i]=new StructMember("currentTime",false,(short)-1,false,(TypeCode)new TypeCode(TCKind.TK_STRING,255),2,false); i++;
        sm[i]=new StructMember("cpuUsage",false,(short)-1,false,(TypeCode)TypeCode.TC_FLOAT,3,false); i++;
        sm[i]=new StructMember("memUsage",false,(short)-1,false,(TypeCode)TypeCode.TC_FLOAT,4,false); i++;
        sm[i]=new StructMember("processes",false,(short)-1,false,(TypeCode)TypeCode.TC_LONG,5,false); i++;

        tc = TypeCodeFactory.TheTypeCodeFactory.create_struct_tc("TopFunction",ExtensibilityKind.EXTENSIBLE_EXTENSIBILITY,sm);
        return tc;
    }
}
